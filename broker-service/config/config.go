package config

type Config struct {
	WebPort string
}

func NewConfig(webPort string) *Config {
	return &Config{
		WebPort: webPort,
	}
}
