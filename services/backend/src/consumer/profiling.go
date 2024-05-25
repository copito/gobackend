package consumer

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"log/slog"
// 	"time"

// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// 	"github.com/copito/data_quality/src/entities"
// 	"github.com/copito/data_quality/src/gateway"
// 	"github.com/spf13/viper"
// )

// func CreateProfileConsumer(ctx context.Context) {
// 	kgateway := gateway.NewKafkaGateway(ctx)

// 	// Consume logic
// 	profileTopic := viper.GetString("kafka.profile_metric_topic")
// 	topics := []string{profileTopic}

// 	reduxFunc := func(ctx context.Context, c *kafka.Consumer, msg *kafka.Message) {
// 		var profileEvent entities.ProfileEvent
// 		err := json.Unmarshal(msg.Value, &profileEvent)
// 		if err != nil {
// 			kgateway.Logger.Warn("unable to unmarshal profile topic", slog.String("err", err.Error()))
// 			return
// 		}

// 		// TODO: Connect to influxdb
// 		var influxClient *gateway.InfluxdbGateway
// 		influxClient = gateway.NewInfluxdbGateway(ctx)
// 		if influxClient != nil {
// 			kgateway.Logger.Warn("unable to create influxdb client")
// 			return
// 		}

// 		// TODO: Push data to timeseries (for profileEvent)
// 		org := viper.GetString("timeseries_db.org")
// 		bucket := viper.GetString("timeseries_db.bucket")
// 		writeAPI := influxClient.Client.WriteAPI(org, bucket)
// 		// writeAPI := influxClient.Client.WriteAPIBlocking(org, bucket)

// 		// p := influxClient.Client.NewPointWithMeasurement("thermostat").
// 		// 	AddTag("unit", "temperature").
// 		// 	AddTag("db_host", profileEvent.DatabaseHost).
// 		// 	AddTag("db_name", profileEvent.DatabaseName).
// 		// 	AddField("min", profileEvent.Payload).
// 		// 	AddField("max", profileEvent.Payload).
// 		// 	SetTime(time.Now())

// 		tags := map[string]string{
// 			"tagname1": "tagvalue1",
// 		}
// 		fields := map[string]interface{}{
// 			"field1": 123,
// 		}
// 		point := writeAPI.NewPoint("measurement1", tags, fields, time.Now())

// 		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
// 			log.Fatal(err)
// 		}
// 		// Flush writes
// 		writeAPI.Flush()
// 	}

// 	kgateway.ConsumeKafkaProfile(ctx, topics, reduxFunc)
// }
