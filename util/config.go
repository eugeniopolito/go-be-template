package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource            string        `mapstructure:"DB_SOURCE"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RedisAddress        string        `mapstructure:"REDIS_ADDRESS"`
	SmtpAuthAddress     string        `mapstructure:"SMTP_AUTH_ADDRESS"`
	SmtpServerAddress   string        `mapstructure:"SMTP_SERVER_ADDRESS"`
	EmailSenderName     string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	VerifyEmailAddress  string        `mapstructure:"VERIFY_EMAIL_ADDRESS"`
	VerifyEmailSubject  string        `mapstructure:"VERIFY_EMAIL_SUBJECT"`
	VerifyEmailBody     string        `mapstructure:"VERIFY_EMAIL_BODY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
