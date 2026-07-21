// Package config loads all runtime configuration from environment
// variables so the binary can be configured without code changes.
package config

import "os"

type Config struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
	// Env is "dev" or "prod". In dev mode templates and static files
	// are read from disk on every request (live reload); in prod they
	// are served from the embedded filesystem.
	Env string
	// DBPath is the path to the SQLite database file.
	DBPath string
	// DefaultLang is the fallback language code, e.g. "en".
	DefaultLang string
	// Dev is true when Env == "dev".
	Dev bool
}

func Load() *Config {
	cfg := &Config{
		Addr:        getenv("ADDR", ":8080"),
		Env:         getenv("ENV", "dev"),
		DBPath:      getenv("DB_PATH", "data/app.db"),
		DefaultLang: getenv("DEFAULT_LANG", "en"),
	}
	cfg.Dev = cfg.Env == "dev"
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
