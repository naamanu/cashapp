package core

import (
	"log"

	"github.com/spf13/viper"
)

type Environment string

const (
	Development Environment = "dev"
	Staging     Environment = "staging"
	Production  Environment = "prod"
)

type Config struct {
	PG_HOST        string `mapstructure:"PG_HOST"`
	PG_PORT        string `mapstructure:"PG_PORT"`
	PG_NAME        string `mapstructure:"PG_NAME"`
	PG_USER        string `mapstructure:"PG_USER"`
	PG_PASS        string `mapstructure:"PG_PASS"`
	PG_SSLMODE     string `mapstructure:"PG_SSLMODE"`
	REDIS_ADDRESS  string `mapstructure:"REDIS_ADDRESS"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	REDIS_DB       int    `mapstructure:"REDIS_DB"`
	REDIS_URL      string `mapstructure:"REDIS_URL"`
	DATABASE_URL   string `mapstructure:"DATABASE_URL"`
	PORT           int    `mapstructure:"PORT"`
	RUN_SEEDS      bool   `mapstructure:"RUN_SEEDS"`
	ENVIRONMENT    Environment
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set Defaults
	viper.SetDefault("PG_HOST", "localhost")
	viper.SetDefault("PG_PORT", "5432")
	viper.SetDefault("PG_NAME", "cashapp")
	viper.SetDefault("PG_USER", "user")
	viper.SetDefault("PG_PASS", "password")
	viper.SetDefault("PG_SSLMODE", "disable")
	viper.SetDefault("REDIS_ADDRESS", "localhost:6379")
	viper.SetDefault("REDIS_DB", 1)
	viper.SetDefault("PORT", 5454)
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("RUN_SEEDS", true)

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we might be using ENV vars
		log.Printf("No .env file found, finding env vars")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	envStr := viper.GetString("ENV")
	if envStr == "" {
		config.ENVIRONMENT = Development
	} else {
		config.ENVIRONMENT = Environment(envStr)
	}

	return &config
}
