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
	}

	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	JWT struct {
		AccessTTL  time.Duration `mapstructure:"access_ttl"`
		RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
	} `mapstructure:"jwt"`

	Hash struct {
		Salt string `mapstructure:"HASH_SALT"`
	}

	Auth struct {
		Secret string `mapstructure:"AUTH_SECRET"`
	}
}

func New(folder, filename, envfilename string) (*Config, error) {
	cfg := new(Config)

	v := viper.New()

	v.SetConfigFile(envfilename + ".env")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg.DB); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg.Hash); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg.Auth); err != nil {
		return nil, err
	}

	v.SetConfigName(filename)
	v.AddConfigPath(folder)
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// if err := v.Unmarshal(&cfg.JWT); err != nil {
	// 	return nil, err
	// }

	return cfg, nil
}
