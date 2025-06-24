package api

import (
	"bitwise74/url-shortener/db"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetURL redirects the user to an origin URL if it exists
func (a *AppRouter) getURL(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No ID provided",
		})
		return
	}

	var record db.ShortURL
	if err := a.DB.First(&record, "short = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No URL found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if record.Origin == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Origin URL is missing",
		})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, record.Origin)
}
