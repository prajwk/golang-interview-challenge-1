package main

func main () {
    // connect to kafka
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

    // read from file
    // publish to kafka

}