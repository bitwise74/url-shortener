// Package api contains everything related to the endpoints and their functions
package api

import (
	"bitwise74/url-shortener/db"
	"bitwise74/url-shortener/logger"
	"bitwise74/url-shortener/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppRouter struct {
	DB  *gorm.DB
	Rt  *gin.Engine
	Log *zap.Logger
}

func SetupApp() (*AppRouter, error) {
	d, err := db.Init()
	if err != nil {
		return nil, err
	}

	r := gin.Default()

	l, err := logger.Init(viper.GetBool("dev"))
	if err != nil {
		return nil, err
	}

	a := &AppRouter{
		DB:  d,
		Rt:  r,
		Log: l,
	}

	if p := viper.GetStringSlice("allowed_proxies"); len(p) == 0 {
		r.SetTrustedProxies(nil)
	} else {
		r.SetTrustedProxies(p)
	}

	// Attach middleware
	r.Use(middleware.RateLimiter())

	// Setup endpoints
	r.POST("/", a.postURL)
	r.GET("/:id", a.getURL)

	if cleanupRate := viper.GetInt("cleanup_interval"); cleanupRate > 0 {
		startDBTicker(d, l, cleanupRate)
	}

	return &AppRouter{
		DB: d,
		Rt: r,
	}, nil
}

func startDBTicker(d *gorm.DB, log *zap.Logger, interval int) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	go func() {
		for range ticker.C {
			err := d.Where("expires_at < ?", time.Now()).Delete(db.ShortURL{}).Error
			if err != nil {
				log.Error("Failed to cleanup database", zap.Error(err))
			}
		}
	}()
}
