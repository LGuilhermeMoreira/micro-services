package config

import (
	"authentication/data"
	"database/sql"
)

type Config struct {
	DB      *sql.DB
	Webport string
	Models  data.Models
}

func NewConfig(db *sql.DB, webport string, models data.Models) *Config {
	return &Config{
		DB:      db,
		Webport: webport,
		Models:  models,
	}
}
