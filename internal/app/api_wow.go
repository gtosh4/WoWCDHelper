package app

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/wow"
)

func registerWoWApi(s *Server) {
	r := s.router.Group("wow")

	r.GET("/classes", s.handleClasses)
	r.GET("/class/:class", s.handleClassInfo)
	r.GET("/class/:class/icon", s.handleClassIcon)
	r.GET("/class/:class/:spec/icon", s.handleClassSpecIcon)

	r.GET("/spec/:specId", s.handleSpecInfo)
	r.GET("/spec/:specId/icon", s.handleSpecIcon)

	r.GET("/spell/:spellId/icon", s.handleSpellIcon)
}

func (s *Server) handleClasses(c *gin.Context) {
	ci, err := wow.ClassInfo(c, s.clients)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ci)
}

func (s *Server) handleClassInfo(c *gin.Context) {
	class := c.Param("class")

	id, err := wow.ClassNameToID(c, s.clients, class)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if id == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cls, _, err := s.clients.Blizz.WoWPlayableClass(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, cls)
}

func (s *Server) handleClassIcon(c *gin.Context) {
	class := c.Param("class")

	id, err := wow.ClassNameToID(c, s.clients, class)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if id == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	s.log.Sugar().Debugf("Found id %d for class %s", id, class)

	media, _, err := s.clients.Blizz.WoWPlayableClassMedia(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.iconResponse(c, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) handleSpecInfo(c *gin.Context) {
	specId, err := strconv.ParseInt(c.Param("specId"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	spec, _, err := s.clients.Blizz.WoWPlayableSpecialization(c, int(specId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, spec)
}

func (s *Server) handleSpecIcon(c *gin.Context) {
	specId, err := strconv.ParseInt(c.Param("specId"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if specId == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	media, _, err := s.clients.Blizz.WoWPlayableSpecializationMedia(c, int(specId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.iconResponse(c, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) handleClassSpecIcon(c *gin.Context) {
	class := c.Param("class")
	spec := c.Param("spec")
	specID, err := wow.ClassSpecToID(c, s.clients, class, spec)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if specID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	media, _, err := s.clients.Blizz.WoWPlayableSpecializationMedia(c, specID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.iconResponse(c, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) handleSpellIcon(c *gin.Context) {
	spell := c.Param("spellId")
	spellId, err := strconv.ParseInt(spell, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	media, _, err := s.clients.Blizz.WoWSpellMedia(c, int(spellId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, asset := range media.Assets {
		if asset.Key == "icon" {
			s.iconResponse(c, asset.Value)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) iconResponse(c *gin.Context, loc string) {
	resp, err := s.clients.IconClient.Get(loc)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
