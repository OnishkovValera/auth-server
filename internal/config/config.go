package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUrl  string `yaml:"DBUrl"`
	JwtKey string `yaml:"jwtKey"`
	Port   string `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
		return nil, err
	}
	return &cfg, nil
}
