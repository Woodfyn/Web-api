package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB     DB
	Server Server
	GRPC   GRPC
	JWT    JWT
	Hash   Hash
	Auth   Auth
}

type DB struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
	Password string `mapstructure:"DB_PASSWORD"`
}

type Server struct {
	Port string
}

type GRPC struct {
	Port string
}

type JWT struct {
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

type Hash struct {
	Salt string `mapstructure:"HASH_SALT"`
}

type Auth struct {
	Secret string `mapstructure:"AUTH_SECRET"`
}

func New(folder, filename, envfilename string) (*Config, error) {
	cfg := &Config{}

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

	v.AddConfigPath(folder)
	v.SetConfigName(filename)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
