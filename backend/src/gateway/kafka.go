package gateway

import "context"

func PublishResultToKafka(ctx context.Context, data []map[string]interface{}) {
	// p, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": viper.GetString("kafka.server"),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer p.Close()

	// topic := viper.GetString("kafka.profile_metric_topic")
	// result := struct {
	// 	data []map[string]interface{}
	// }{
	// 	data: profileResult,
	// }
	// value, err := json.Marshal(result)
	// if err != nil {
	// 	panic(err)
	// }

	// // Publish data to kafka topic to be consumed by the consumer (that sends to influxdb)
	// err = p.Produce(&kafka.Message{
	// 	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 	Value:          value,
	// }, nil)
	// if err != nil {
	// 	panic(err)
	// }
}
