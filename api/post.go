package api

import (
	"bitwise74/url-shortener/db"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const charsetLen = len(charset)

type urlInput struct {
	URL string `form:"url"`
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func makeShortURL(n int) string {
	if n <= 0 {
		return ""
	}

	var s strings.Builder
	s.Grow(n)

	for range n {
		s.WriteByte(charset[rand.Intn(charsetLen)])
	}

	return s.String()
}

// postURL takes the provided URL and creates a new short URL that's returned the user
func (a *AppRouter) postURL(c *gin.Context) {
	var input urlInput

	if err := c.ShouldBind(&input); err != nil {
		fmt.Println(err)
	}

	if input.URL == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "URL field can't be empty",
		})
		return
	}

	if !isURL(input.URL) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Provided URL isn't valid",
		})
		return
	}

	record := db.ShortURL{}

	if err := a.DB.First(&record, "origin = ?", input.URL).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}
	}

	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}

	hostURL := scheme + c.Request.Host + "/"

	// Check if we have a record in the database already
	if record.Origin != "" && record.Short != "" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"short_url": hostURL + record.Short,
		})
		return
	}

	// If not make a new record and send the response
	short := makeShortURL(viper.GetInt("url_id_size"))
	expiresAt := time.Now().Add(time.Second * 60).Unix()

	if err := a.DB.Create(&db.ShortURL{Origin: input.URL, Short: short, ExpiresAt: expiresAt}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"short_url": hostURL + short,
	})
}
