package config

type Config struct {
	WebPort string
}

func NewConfig(webport string) *Config {
	return &Config{
		WebPort: webport,
	}
}
