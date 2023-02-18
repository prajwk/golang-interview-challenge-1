// package main
//
// func main () {
//     // connect to kafka
//     kafkaReader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:     cfg.Kafka.Brokers,
// 		GroupID:     cfg.Kafka.GroupID,
// 		GroupTopics: []string{cfg.Kafka.Topics.IntradayValues},
// 		ErrorLogger: kafka.LoggerFunc(log.Errorf),
// 	})
// 	defer kafkaReader.Close()
//
// 	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry.URL))
// 	if err != nil {
// 		log.WithError(err).Fatal("not connected to schema registry")
// 	}
//
//     // read from file
//     // publish to kafka
//
// }


package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/confluentinc/confluent-kafka-go/schemaregistry"
    "github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
)

func main() {
    // Define Kafka configuration
    kafkaConf := &kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "acks":              "all",
    }

    // Create Kafka producer
    producer, err := kafka.NewProducer(kafkaConf)
    if err != nil {
        log.Fatalf("Failed to create producer: %s", err)
    }
    defer producer.Close()

    // Define schema registry configuration
    schemaRegistryConf := &schemaregistry.Config{
        URL: "http://localhost:8081",
    }

    // Create schema registry client
    schemaRegistryClient, err := schemaregistry.NewClient(schemaRegistryConf)
    if err != nil {
        log.Fatalf("Failed to create schema registry client: %s", err)
    }
    defer schemaRegistryClient.Close()

    // Read JSON file
    jsonFile, err := ioutil.ReadFile("test/testdata/intraday-values.json")
    if err != nil {
        log.Fatalf("Failed to read JSON file: %s", err)
    }

    // Parse JSON
    var data interface{}
    err = json.Unmarshal(jsonFile, &data)
    if err != nil {
        log.Fatalf("Failed to parse JSON: %s", err)
    }

    // Get schema ID from schema registry
    schemaID, err := schemaRegistryClient.GetLatestSchemaID("example-value")
    if err != nil {
        log.Fatalf("Failed to get schema ID: %s", err)
    }

    // Serialize data using Avro schema
    avroBytes, err := serde.AvroFromNative(schemaRegistryClient, schemaID, data)
    if err != nil {
        log.Fatalf("Failed to serialize data: %s", err)
    }

    // Publish data to Kafka topic
    err = producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &topicName,
            Partition: kafka.PartitionAny,
        },
        Key:   nil,
        Value: avroBytes,
        Headers: []kafka.Header{
            kafka.Header{
                Key:   "schemaID",
                Value: []byte(fmt.Sprintf("%d", schemaID)),
            },
        },
    }, nil)

    if err != nil {
        log.Fatalf("Failed to publish message: %s", err)
    }

    log.Printf("Message published successfully to topic %s", topicName)
}
