package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Name                 string
	Port                 int
	Env                  string
	LogLevel             string `mapstructure:"log_level"`
	MarketUpdateInterval string `mapstructure:"market_update_interval"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int `mapstructure:"expiration_hours"`
}

var App *Config

func LoadConfig(dir string) {
	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	env := viper.GetString("app.env")
	if env != "" && env != "config" {
		viper.SetConfigName(env)
		if err := viper.MergeInConfig(); err != nil {
			log.Printf("No environment specific config found for %s, or error: %s", env, err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	App = &config
	log.Println("Configuration loaded successfully. Environment:", App.App.Env)
}
