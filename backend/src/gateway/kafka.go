package gateway

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

func PublishResultToKafka(ctx context.Context, topic string, key string, data string) {
	logger := ctx.Value("logger").(*slog.Logger)
	kafkaBrokerHost := viper.GetString("kafka.server")
	if kafkaBrokerHost == "" {
		logger.Error("kafka brokers not provided to push profiling data", slog.String("servers", kafkaBrokerHost))
	}

	if topic == "" {
		logger.Error("kafka topic not provided to push profiling data", slog.String("topic", topic))
	}

	// Connect to kafka host
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerHost,
	})
	if err != nil {
		logger.Error("unable to connect to kafka broker", slog.String("servers", kafkaBrokerHost), slog.String("err", err.Error()))
		panic(err)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func(p *kafka.Producer) {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Error("delivery failed", slog.Any("topic_partition", ev.TopicPartition))
					// rollback strategy applying
				} else {
					logger.Info("delivery success", slog.Any("topic_partition", ev.TopicPartition))
				}
			}
		}
	}(p)

	// Publish data to kafka topic to be consumed by the consumer (that sends to influxdb)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
		Key:            []byte(key),
	}, nil)
	if err != nil {
		logger.Error(
			"unable to publish to kafka topic",
			slog.String("servers", kafkaBrokerHost),
			slog.String("topic", topic),
			slog.String("key", key),
			slog.Any("result", data),
		)
		panic(err)
	}

	p.Flush(15 * 1000)
}
