package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/wow"
)

func registerWoWApi(s *Server) {
	r := s.router.Group("wow")

	r.GET("/class-info", s.handleClassInfo)
	r.GET("/class/:class/icon", s.handleClassIcon)
	r.GET("/class/:class/:spec/icon", s.handleSpecIcon)

	r.GET("/spell/:spellId/icon", s.handleSpellIcon)
}

func (s *Server) handleClassInfo(c *gin.Context) {
	ci, err := wow.ClassInfo(c.Request.Context(), s.clients)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ci)
}

func (s *Server) handleClassIcon(c *gin.Context) {
	class := c.Param("class")

	id, err := wow.ClassNameToID(c.Request.Context(), s.clients, class)
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

func (s *Server) handleSpecIcon(c *gin.Context) {
	class := c.Param("class")
	spec := c.Param("spec")
	specID, err := wow.ClassSpecToID(c.Request.Context(), s.clients, class, spec)
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

func (s *Server) handleSpellIcon(c *gin.Context) {
	spell := c.Param("spellId")
	spellId, err := strconv.ParseInt(spell, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	media, _, err := s.clients.Blizz.WoWSpellMedia(c.Request.Context(), int(spellId))
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