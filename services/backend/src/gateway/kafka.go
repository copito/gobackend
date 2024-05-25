package gateway

import (
	"context"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

type KafkaGateway struct {
	Logger   *slog.Logger
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

type ReduxFunc func(ctx context.Context, c *kafka.Consumer, msg *kafka.Message)

func NewKafkaGateway(ctx context.Context) *KafkaGateway {
	logger := ctx.Value("logger").(*slog.Logger)

	kafkaBrokerHost := viper.GetString("kafka.server")
	if kafkaBrokerHost == "" {
		logger.Error("kafka brokers not provided to push profiling data", slog.String("servers", kafkaBrokerHost))
	}

	// Connect to kafka host
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerHost,
	})
	if err != nil {
		logger.Error("unable to connect to kafka broker", slog.String("servers", kafkaBrokerHost), slog.String("err", err.Error()))
		panic(err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerHost,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	return &KafkaGateway{
		Logger:   logger,
		Producer: p,
		Consumer: c,
	}
}

func (k *KafkaGateway) PublishResultToKafka(ctx context.Context, topic string, key string, data string) {
	logger := ctx.Value("logger").(*slog.Logger)

	if topic == "" {
		logger.Error("kafka topic not provided to push profiling data", slog.String("topic", topic))
	}

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
	}(k.Producer)

	// Publish data to kafka topic to be consumed by the consumer (that sends to influxdb)
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
		Key:            []byte(key),
	}, nil)
	// Error checking publish
	if err != nil {
		logger.Error(
			"unable to publish to kafka topic",
			slog.String("topic", topic),
			slog.String("key", key),
			slog.Any("result", data),
		)
		panic(err)
	}

	k.Producer.Flush(15 * 1000)
}

func (k *KafkaGateway) ConsumeKafkaProfile(ctx context.Context, topics []string, redux ReduxFunc) {
	if len(topics) == 0 {
		k.Logger.Error("kafka topic not provided to push profiling data")
		panic("no topics provided")
	}

	// c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)
	err := k.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		k.Logger.Error("unable to subcribe to topics", slog.Any("topic", topics))
		panic("unable to subcribe to topics")
	}

	// A signal handler or similar could be used to set this to false to break the loop.
	for {
		msg, err := k.Consumer.ReadMessage(3 * time.Second)
		if err != nil {
			k.Logger.Error(
				"Consumer error",
				slog.String("err", err.Error()),
				slog.String("msg", msg.String()),
			)
		}

		if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			k.Logger.Error(
				"Timeout Consumer error",
				slog.String("err", err.Error()),
				slog.String("msg", msg.String()),
			)
			break
		}

		k.Logger.Info(
			"Message on topic partition",
			slog.String("topic", msg.TopicPartition.String()),
			slog.String("msg", string(msg.Value)),
		)

		// Perform redux action
		redux(ctx, k.Consumer, msg)

	}
}

func (k *KafkaGateway) Close() error {
	err := k.Consumer.Close()
	if err != nil {
		k.Logger.Error("error closing kafka consumer", slog.String("err", err.Error()))
		return err
	}

	k.Producer.Close()
	return nil
}
