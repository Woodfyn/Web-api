package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     string `mapstructure:"DB_PORT"`
		Username string `mapstructure:"DB_USERNAME"`
		Name     string `mapstructure:"DB_NAME"`
		SSLMode  string `mapstructure:"DB_SSLMODE"`
		Password string `mapstructure:"DB_PASSWORD"`
	} `mapstructure:"db"`

	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Hash struct {
		Salt string `mapstructure:"HASH_SALT"`
	} `mapstructure:"hash"`

	JWT struct {
		TokenTTL time.Duration `mapstructure:"token_ttl"`
	} `mapstructure:"auth"`

	Auth struct {
		Secret string `mapstructure:"AUTH_SECRET"`
	} `mapstructure:"auth"`
}

func New(envfilename, folder, filename string) (*Config, error) {
	cfg := new(Config)
	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	viper.SetConfigFile(envfilename + ".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg.DB); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.Hash); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.Auth); err != nil {
		return nil, err
	}

	return cfg, nil
}
