package app

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/FuzzyStatic/blizzard/v2/wowgd"
	"github.com/gin-gonic/gin"
)

func registerApi(s *Server) {
	r := s.router.Group("api")

	r.GET("/class-info", s.handleClassInfo)
	r.GET("/class/:class/icon", s.handleClassIcon)
	r.GET("/class/:class/:spec/icon", s.handleSpecIcon)
}

func (s *Server) classInfo(ctx context.Context) (classes []wowgd.PlayableClass, err error) {
	idx, _, err := s.clients.Blizz.WoWPlayableClassesIndex(ctx)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	classes = make([]wowgd.PlayableClass, len(idx.Classes))
	for i, cl := range idx.Classes {
		wg.Add(1)
		go func(i int, id int) {
			defer wg.Done()
			c, _, err := s.clients.Blizz.WoWPlayableClass(ctx, id)
			if err != nil {
				s.log.Sugar().Warnf("error loading class %d: %v", id, err)
			} else {
				classes[i] = *c
			}
		}(i, cl.ID)
	}
	wg.Wait()

	return classes, nil
}

func (s *Server) handleClassInfo(c *gin.Context) {
	ci, err := s.classInfo(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ci)
}

func (s *Server) classNameToID(ctx context.Context, name string) (int, error) {
	var id int
	ci, err := s.classInfo(ctx)
	if err != nil {
		return 0, err
	}
	for _, info := range ci {
		if strings.EqualFold(name, info.Name) {
			id = info.ID
			break
		}
	}
	if id == 0 {
		i64, err := strconv.ParseInt(name, 10, 64)
		if err == nil {
			id = int(i64)
		}
	}

	return id, nil
}

func (s *Server) handleClassIcon(c *gin.Context) {
	class := c.Param("class")

	id, err := s.classNameToID(c.Request.Context(), class)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if id == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	s.log.Sugar().Debugf("Found id %d for class %s", id, class)

	media, _, err := s.clients.Blizz.WoWPlayableClassMedia(c.Request.Context(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.log.Sugar().Debugf("redirecting to %s", asset.Value)
			c.Redirect(http.StatusTemporaryRedirect, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) classSpecToID(ctx context.Context, class string, spec string) (int, error) {
	ci, err := s.classInfo(ctx)
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

func (s *Server) handleSpecIcon(c *gin.Context) {
	class := c.Param("class")
	spec := c.Param("spec")
	specID, err := s.classSpecToID(c.Request.Context(), class, spec)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if specID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	media, _, err := s.clients.Blizz.WoWPlayableSpecializationMedia(c.Request.Context(), specID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.log.Sugar().Debugf("redirecting to %s", asset.Value)
			c.Redirect(http.StatusTemporaryRedirect, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}
