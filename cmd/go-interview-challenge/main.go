package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gpsinsight/go-interview-challenge/internal/config"
	"github.com/gpsinsight/go-interview-challenge/internal/consumer"
	"github.com/gpsinsight/go-interview-challenge/internal/server"
	"github.com/gpsinsight/go-interview-challenge/pkg/messages"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	log := logger.WithField("environment", "local")

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config.New", err)
	}

	// Setup context that will cancel on signalled termination
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		log.Info("termination signaled")
		cancel()
	}()

	db := sqlx.MustConnect("postgres", cfg.Postgres.ConnectionString())
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(time.Hour)
	defer db.Close()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Kafka.Brokers,
		GroupID:     cfg.Kafka.GroupID,
		GroupTopics: []string{cfg.Kafka.Topics.IntradayValues},
		ErrorLogger: kafka.LoggerFunc(log.Errorf),
	})
	defer kafkaReader.Close()

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry.URL))
	if err != nil {
		log.WithError(err).Fatal("not connected to schema registry")
	}

	protoDeserializer, err := protobuf.NewDeserializer(client, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		log.WithError(err).Fatal("failed to set up deserializer")
	}
	err = protoDeserializer.ProtoRegistry.RegisterMessage((&messages.IntradayValue{}).ProtoReflect().Type())
	if err != nil {
		log.WithError(err).Fatal("failed to register IntradayValue message type")
	}

	kafkaConsumer := consumer.NewKafkaConsumer(
		kafkaReader,
		protoDeserializer,
		db,
		log,
	)
	go kafkaConsumer.Run(ctx)

	log.Info("Starting up go-interview-challenge")

	srvr := server.New(cfg, log)
	defer func() {
		err := srvr.Shutdown(context.Background())
		if err != nil {
			logger.WithError(err).Error("failed to gracefully shutdown server")
		}
	}()
	go func() {
		err := srvr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	// Exit safely
	<-ctx.Done()
	log.Info("exiting")
}
