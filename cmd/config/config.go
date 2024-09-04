package config

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	customerrors "github.com/simondanielsson/apPRoved/pkg/custom_errors"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	BindAddr string `mapstructure:"bind_address"`
	Mode     string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	JWT      struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	envFile := ".env." + env

	log.Printf("loading config from %s\n", envFile)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	viper.SetConfigFile(path.Join("config", "config.yaml"))
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	customerrors.IgnoreError(viper.BindEnv("database.user", "POSTGRES_USER"))
	customerrors.IgnoreError(viper.BindEnv("database.password", "POSTGRES_PASSWORD"))
	customerrors.IgnoreError(viper.BindEnv("database.host", "POSTGRES_HOST"))
	customerrors.IgnoreError(viper.BindEnv("database.port", "POSTGRES_PORT"))
	customerrors.IgnoreError(viper.BindEnv("database.dbname", "POSTGRES_DBNAME"))
	customerrors.IgnoreError(viper.BindEnv("jwt.secret", "JWT_KEY"))

	// discover the config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Server == nil {
		log.Fatalf("server config is missing")
	}
	if cfg.Database == nil {
		log.Fatalf("database config is missing")
	}

	return &cfg, nil
}
