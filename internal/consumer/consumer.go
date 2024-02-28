package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"github.com/gpsinsight/go-interview-challenge/pkg/messages"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	reader       *kafka.Reader
	deserializer *protobuf.Deserializer
	db           *sqlx.DB
	logger       *logrus.Entry
}

func NewKafkaConsumer(
	reader *kafka.Reader,
	deserializer *protobuf.Deserializer,
	db *sqlx.DB,
	logger *logrus.Entry,
) *KafkaConsumer {
	return &KafkaConsumer{
		reader:       reader,
		deserializer: deserializer,
		db:           db,
		logger:       logger,
	}
}

func (kc *KafkaConsumer) Run(ctx context.Context) {
	for {
		msg, err := kc.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				kc.logger.Info("Context cancelled. Stopping consumer...")
				break
			}
			kc.logger.WithError(err).Error("Failed to read from kafka")
			continue
		}

		err = kc.processMessage(ctx, msg)
		if err != nil {
			kc.logger.WithError(err).Error("Failed to process message")
			continue
		}

		err = kc.reader.CommitMessages(ctx, msg)
		if err != nil {
			kc.logger.WithError(err).Error("Failed to commit message")
		}
	}
}

func (kc *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) error {
	// Deserialize the Protobuf message
	intradayVal, err := kc.deserializer.Deserialize(kc.reader.Config().Topic, msg.Value)
	if err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}
	data, ok := intradayVal.(*messages.IntradayValue)
	if !ok {
		return fmt.Errorf("failed to deserialize message: not of type IntradayValue")
	}
	// Insert the deserialized data into the PostgreSQL table
	query := `INSERT INTO intraday (ticker, timestamp, open, high, low, close, volume) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = kc.db.ExecContext(ctx, query, data.Ticker, data.Timestamp, data.Open, data.High, data.Low, data.Close, data.Volume)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %w", err)
	}
	kc.logger.Infof("received message: %s", string(msg.Key))
	kc.logger.Infof("Inserted data: ticker=%s, timestamp=%s, open=%f, high=%f, low=%f, close=%f, volume=%d into Postgresql", data.Ticker, data.Timestamp, data.Open, data.High, data.Low, data.Close, data.Volume)
	return nil
}
