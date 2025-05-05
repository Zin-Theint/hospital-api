package config

import "os"

type Config struct {
	Port      string
	DB_DSN    string
	JWTSecret string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev‑secret‑change‑me"
	}
	return Config{
		Port:      port,
		DB_DSN:    os.Getenv("DB_DSN"),
		JWTSecret: secret,
	}
}
