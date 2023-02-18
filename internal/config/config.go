package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type TopicConfig struct {
	IntradayValues string `envconfig:"KAFKA_INTRADAY_VALUE" default:"gpsi.stocks.intraday.value"`
}

type SchemaRegistryConfig struct {
	URL                 string `envconfig:"KAFKA_REGISTRY_URL" default:"http://localhost:8081"`
	Username            string `envconfig:"KAFKA_REGISTRY_USER"`
	Password            string `envconfig:"KAFKA_REGISTRY_PASSWORD"`
	AutoRegisterSchemas bool   `envconfig:"KAFKA_REGISTRY_AUTO_REGISTER" default:"false"`
}

type KafkaConfig[T any] struct {
	Brokers                []string `envconfig:"KAFKA_BROKER_HOSTS" default:"localhost:9092"`
	GroupID                string   `envconfig:"KAFKA_CONSUMER_GROUP_ID"`
	Username               string   `envconfig:"KAFKA_USER"`
	Password               string   `envconfig:"KAFKA_PASSWORD"`
	AllowAutoTopicCreation bool     `envconfig:"KAFKA_AUTO_TOPIC_CREATION" default:"false"`
	Topics                 T
	SchemaRegistry         SchemaRegistryConfig
}

type PostgresConfig struct {
	Password     string `envconfig:"POSTGRES_PASSWORD"`
	User         string `envconfig:"POSTGRES_USER"     default:"postgres"`
	Port         int    `envconfig:"POSTGRES_PORT"     default:"5432"`
	Database     string `envconfig:"POSTGRES_DB"       default:"postgres"`
	Host         string `envconfig:"POSTGRES_HOST"`
	SSLMode      string `envconfig:"POSTGRES_SSL_MODE" default:"require"`
	MaxOpenConns int    `envconfig:"POSTGRES_MAX_OPEN_CONNS"`
}

// ConnectionString builds a postgres connection string from the configured values
func (pc PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.Host,
		pc.Port,
		pc.User,
		pc.Password,
		pc.Database,
		pc.SSLMode,
	)
}

type Config struct {
	Environment string `envconfig:"ENVIRONMENT"`
	Postgres    PostgresConfig
	Kafka       KafkaConfig[TopicConfig]
}

// New creates a new Config from environment variables and env files in a specified directory
func New() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
