package warcraftlogs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/avast/retry-go"
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/ctx"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/fight"
)

type Client struct {
	*ctx.Ctx

	HTTP *http.Client
	Base *url.URL
	DB   *badger.DB

	limiter *rate.Limiter
}

const ApiKeyParam = "api_key"

var retryRequestErr = errors.New("Retryable request error")

func NewClient(ctx *ctx.Ctx, HTTP *http.Client, baseUrl string, apiKey string, db *badger.DB) *Client {
	u, _ := url.Parse(baseUrl)
	q := u.Query()
	q.Set(ApiKeyParam, apiKey)
	q.Set("translate", "true")
	u.RawQuery = q.Encode()

	ctx.Log.Infof("Using base WarcraftLogs URL of %s", u)

	c := &Client{
		Ctx:  ctx,
		HTTP: HTTP,
		Base: u,
		DB:   db,

		limiter: rate.NewLimiter(rate.Limit(3), 30),
	}

	return c
}

func (c *Client) get(path string, params url.Values, out interface{}) (err error) {
	u, _ := url.Parse(path)
	u.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out interface{}) (err error) {
	log := c.Log.WithField("method", req.Method).WithField("url", req.URL.String())

	u := *req.URL
	u.Scheme = c.Base.Scheme
	u.Host = c.Base.Host
	u.Path = strings.TrimSuffix(c.Base.Path, "/") + "/" + strings.TrimPrefix(req.URL.Path, "/")
	params := u.Query()
	for k, vs := range c.Base.Query() {
		for _, v := range vs {
			params.Add(k, v)
		}
	}
	u.RawQuery = params.Encode()
	req.URL = &u
	req.Header.Add("Accept-Encoding", "gzip")
	req = req.WithContext(c.Context)

	var resp *http.Response
	buf := new(bytes.Buffer)
	err = retry.Do(
		func() error {
			if err := c.limiter.Wait(c.Context); err != nil {
				return err
			}

			resp, err = c.HTTP.Do(req)
			if err != nil {
				return err
			}
			switch c := resp.StatusCode; {
			case c < 200:
				return errors.Errorf("Request returned with status code %d", resp.StatusCode)
			case c == 429:
				log.Debug("Received status code 429, retrying")
				return retryRequestErr
			case c >= 400:
				return errors.Errorf("Request returned with status code %d", resp.StatusCode)
			}
			defer resp.Body.Close()

			var body io.ReadCloser
			switch resp.Header.Get("Content-Encoding") {
			case "gzip":
				body, err = gzip.NewReader(resp.Body)
				if err != nil {
					return err
				}

			default:
				body = ioutil.NopCloser(resp.Body)
			}
			defer body.Close()

			_, err = io.Copy(buf, body)
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(func(err error) bool {
			if err == retryRequestErr {
				return false
			}
			return false
		}),
	)
	if err != nil {
		return
	}

	switch ct := resp.Header.Get("Content-Type"); {
	case strings.HasPrefix(ct, "application/json"):
		err = errors.Wrap(json.Unmarshal(buf.Bytes(), out), "Error unmarshaling JSON response")

	default:
		err = errors.Errorf("unsupported Content-Type: %s", resp.Header.Get("Content-Type"))
		if err != nil {
			c.Log.Debugf("Bad content type for body:\n%s", buf)
		}
	}
	return
}

func (c *Client) Fights(code string) (fs fight.Fights, err error) {
	key := []byte("fights_" + code)
	log := c.Log.WithField("cache-key", string(key))

	err = c.DB.Update(func(txn *badger.Txn) error {
		item, rerr := txn.Get(key)
		if rerr == nil {
			rerr = item.Value(func(val []byte) error {
				return (&fs).Unmarshal(val)
			})
			if rerr != nil {
				log.WithError(rerr).Warnf("Could not unmarshal cached value")
			} else {
				return nil
			}
		}
		rerr = c.get("/report/fights/"+code, nil, &fs)
		if rerr != nil {
			return rerr
		}

		entry := &badger.Entry{
			Key: key,
		}
		entry.Value, rerr = fs.Marshal()
		if rerr = txn.SetEntry(entry); rerr != nil {
			log.WithError(rerr).Infof("Cache put failed")
		}
		return nil
	})
	return
}
