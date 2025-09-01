package config

import (
	"log"
	"os"
)

func MustLoad() *Config {
	cfg := Config{
		Port:        getenv("PORT", DefaultPort),
		Debug:       getenv("DEBUG", DefaultDebug) == "true",
		StaticToken: getenv("STATIC_TOKEN", DefaultStaticToken),
	}
	log.Printf("cfg: %+v", cfg)
	return &cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
