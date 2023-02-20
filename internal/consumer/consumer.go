package consumer

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
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
	/**
	   * TODO:
		 * - deserialize protobuf message
		 * - insert data into postgres table
	*/

	kc.logger.Infof("received message: %s", string(msg.Key))

	return nil
}
