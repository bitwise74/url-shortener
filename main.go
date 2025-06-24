package main

import (
	"bitwise74/url-shortener/api"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port            int      `mapstructure:"port"`
	Dev             bool     `mapstructure:"dev"`
	AllowedProxies  []string `mapstructure:"allowed_proxies"`
	RateLimiterMode string   `mapstructure:"rate_limiter_mode"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	viper.BindEnv("port", "PORT")
	viper.BindEnv("dev", "DEV")
	viper.BindEnv("allowed_proxies", "ALLOWED_PROXIES")
	viper.BindEnv("rate_limiter_mode", "RATE_LIMITER_MODE")

	// Validate some config options
	if p := viper.GetString("port"); p == "" {
		panic("port is missing")
	}

	if i := viper.GetInt("url_id_size"); i < 2 || i > 200 {
		panic("url_id_size is invalid")
	}

	// SSL stuff here
	if s := viper.GetBool("secure"); s {
		if viper.GetString("ssl_cert_path") == "" {
			panic("missing ssl_cert_path with ssl enabled")
		}

		if viper.GetString("ssl_key_path") == "" {
			panic("missing ssl_key_path with ssl enabled")
		}
	}

	app, err := api.SetupApp()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("secure") {
		app.Rt.RunTLS(fmt.Sprintf(":%v", viper.GetInt("port")), viper.GetString("ssl_cert_path"), viper.GetString("ssl_key_path"))
	}
	app.Rt.Run(fmt.Sprintf(":%v", viper.GetInt("port")))
}
