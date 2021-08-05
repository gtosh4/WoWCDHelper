package wow

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/FuzzyStatic/blizzard/wowgd"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/pkg/errors"
)

func ClassInfo(ctx context.Context, clients *clients.Clients) (classes []wowgd.PlayableClass, err error) {
	idx, _, err := clients.Blizz.WoWPlayableClassesIndex(ctx)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	classes = make([]wowgd.PlayableClass, len(idx.Classes))
	for i, cl := range idx.Classes {
		wg.Add(1)
		go func(i int, id int) {
			defer wg.Done()
			c, _, err := clients.Blizz.WoWPlayableClass(ctx, id)
			if err != nil {
				clients.Log.Sugar().Warnf("error loading class %d: %v", id, err)
			} else {
				classes[i] = *c
			}
		}(i, cl.ID)
	}
	wg.Wait()

	return classes, nil
}

func ClassNameToID(ctx context.Context, clients *clients.Clients, name string) (int, error) {
	id, err := strconv.ParseInt(name, 10, 63)
	if err == nil {
		return int(id), nil
	}
	ci, err := ClassInfo(ctx, clients)
	if err != nil {
		return 0, err
	}
	for _, info := range ci {
		if strings.EqualFold(name, info.Name) {
			return info.ID, nil
		}
	}

	return 0, errors.Errorf("class '%s' not found", name)
}

func ClassSpecToID(ctx context.Context, clients *clients.Clients, class string, spec string) (int, error) {
	ci, err := ClassInfo(ctx, clients)
	if err != nil {
		return 0, err
	}
	for _, classInfo := range ci {
		if strings.EqualFold(class, classInfo.Name) {
			for _, specInfo := range classInfo.Specializations {
				if strings.EqualFold(spec, specInfo.Name) {
					return specInfo.ID, nil
				}
			}
		}
	}

	return 0, nil
}
