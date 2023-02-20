package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"github.com/gpsinsight/go-interview-challenge/internal/config"
	"github.com/gpsinsight/go-interview-challenge/pkg/messages"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type IntradayValue struct {
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

func main() {
	var cfg config.Config
	_ = envconfig.Process("", &cfg)
	log := logrus.NewEntry(logrus.New())

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers...),
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.Hash{},
		Topic:                  cfg.Kafka.Topics.IntradayValues,
	}
	defer w.Close()

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry.URL))
	if err != nil {
		log.WithError(err).Fatal("not connected to schema registry")
	}
	serConfig := protobuf.NewSerializerConfig()
	serConfig.AutoRegisterSchemas = true
	protoSerializer, err := protobuf.NewSerializer(client, serde.ValueSerde, serConfig)
	if err != nil {
		log.WithError(err).Error("failed to set up serializer")
		return
	}

	jsonFile, err := ioutil.ReadFile("test/testdata/intraday-values.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}

	var data map[string]IntradayValue
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	for key, value := range data {
		open, _ := strconv.ParseFloat(value.Open, 64)
		high, _ := strconv.ParseFloat(value.High, 64)
		low, _ := strconv.ParseFloat(value.Low, 64)
		close, _ := strconv.ParseFloat(value.Close, 64)
		volume, _ := strconv.ParseInt(value.Volume, 10, 64)

		kafkaValue := &messages.IntradayValue{
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
			Timestamp: key,
			Ticker:    "IBM",
		}

		payload, err := protoSerializer.Serialize(cfg.Kafka.Topics.IntradayValues, kafkaValue)
		if err != nil {
			log.WithError(err).Fatal("failed to serialize")
		}

		msg := kafka.Message{
			Key:   []byte(key),
			Value: payload,
		}

		w.WriteMessages(context.Background(), msg)

		log.Infof("Successfully published message: %s", key)
	}
}
