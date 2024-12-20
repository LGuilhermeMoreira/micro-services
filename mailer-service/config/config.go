package config

import (
	"mailer-service/email"
	"os"
	"strconv"
)

type Config struct {
	WebPort string
	Mailer  email.Mail
}

func NewConfig(webPort string) *Config {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return &Config{
		WebPort: webPort,
		Mailer: email.Mail{
			Port:        port,
			Domain:      os.Getenv("MAIL_DOMAIN"),
			Host:        os.Getenv("MAIL_HOST"),
			UserName:    os.Getenv("MAIL_USERNAME"),
			Password:    os.Getenv("MAIL_PASSWORD"),
			Encryption:  os.Getenv("MAIL_ENCRYPTION"),
			FromName:    os.Getenv("FROM_NAME"),
			FromAddress: os.Getenv("FROM_ADDRESS"),
		}}
}
