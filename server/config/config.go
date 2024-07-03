// config/config.go
package config

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	Port                               string `mapstructure:"PORT"`
	DbHost                             string `mapstructure:"DB_HOST"`
	DbName                             string `mapstructure:"DB_NAME"`
	DbUser                             string `mapstructure:"DB_USER"`
	DbPassword                         string `mapstructure:"DB_PASSWORD"`
	DbPort                             string `mapstructure:"DB_PORT"`
	DbSslMode                          string `mapstructure:"DB_SSL_MODE"`
	JwtAccessTokenSecret               string `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	JwtAccessTokenExpirationInSeconds  uint   `mapstructure:"JWT_ACCESS_TOKEN_EXPIRATION_SECONDS"`
	JwtRefreshTokenSecret              string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
	JwtRefreshTokenExpirationInSeconds uint   `mapstructure:"JWT_REFRESH_TOKEN_EXPIRATION_SECONDS"`
}

var Env *EnvConfig

func LoadConfig() error {
	viper.SetConfigFile(".env") // required if the config file does not have the extension in the name
	viper.AutomaticEnv()        // automatically override values from environment variables

	// Set default values
	viper.SetDefault("PORT", 3000)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var cfg EnvConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	Env = &cfg

	return nil
}
