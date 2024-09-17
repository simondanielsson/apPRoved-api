package config

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	customerrors "github.com/simondanielsson/apPRoved/pkg/custom_errors"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	BindAddr string `mapstructure:"bind_address"`
	Mode     string `mapstructure:"mode"`
	AMQPMode string `mapstructure:"amqp_mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type RabbitMQQueueConfig struct {
	Name       string     `mapstructure:"name"`
	Durable    bool       `mapstructure:"durable"`
	AutoDelete bool       `mapstructure:"auto_delete"`
	Exclusive  bool       `mapstructure:"exclusive"`
	NoWait     bool       `mapstructure:"no_wait"`
	Args       amqp.Table `mapstructure:"args"`
}

type RabbitMQConfig struct {
	Url    string                `mapstructure:"url"`
	Queues []RabbitMQQueueConfig `mapstructure:"queues"`
}

type PubSubConfig struct {
	ProjectID string `mapstructure:"project_id"`
}

type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	JWT      *JWTConfig      `mapstructure:"jwt"`
	MQ       *RabbitMQConfig `mapstructure:"mq"`
	PubSub   *PubSubConfig   `mapstructure:"pubsub"`
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
	parseEnvironmentVariables()

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
	if cfg.JWT == nil {
		log.Fatalf("jwt config is missing")
	}
	if cfg.MQ == nil {
		log.Fatalf("rabbitmq config is missing")
	}
	if cfg.PubSub == nil {
		log.Fatalf("pubsub config is missing")
	}

	if err := ValidateRabbitMQConfig(cfg.MQ); err != nil {
		log.Fatalf("configuration validation error: %v", err)
	}

	return &cfg, nil
}

func parseEnvironmentVariables() {
	customerrors.IgnoreError(viper.BindEnv("server.mode", "APP_ENV"))
	customerrors.IgnoreError(viper.BindEnv("server.bind_address", "APP_PORT"))
	customerrors.IgnoreError(viper.BindEnv("server.amqp_mode", "AMQP_MODE"))

	customerrors.IgnoreError(viper.BindEnv("database.user", "POSTGRES_USER"))
	customerrors.IgnoreError(viper.BindEnv("database.user", "POSTGRES_USER"))
	customerrors.IgnoreError(viper.BindEnv("database.password", "POSTGRES_PASSWORD"))
	customerrors.IgnoreError(viper.BindEnv("database.host", "POSTGRES_HOST"))
	customerrors.IgnoreError(viper.BindEnv("database.port", "POSTGRES_PORT"))
	customerrors.IgnoreError(viper.BindEnv("database.dbname", "POSTGRES_DBNAME"))

	customerrors.IgnoreError(viper.BindEnv("jwt.secret", "JWT_KEY"))
	customerrors.IgnoreError(viper.BindEnv("mq.url", "AMQP_URL"))
	customerrors.IgnoreError(viper.BindEnv("pubsub.project_id", "GCP_PROJECT_ID"))
}
