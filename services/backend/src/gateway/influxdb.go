package gateway

import (
	"context"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
)

type InfluxdbGateway struct {
	Logger *slog.Logger
	Client *influxdb2.Client
}

func NewInfluxdbGateway(ctx context.Context) *InfluxdbGateway {
	logger := ctx.Value("logger").(*slog.Logger)

	dbToken := viper.GetString("timeseries_db.token")
	if dbToken == "" {
		logger.Error("influxdb token not set")
	}

	dbURL := viper.GetString("timeseries_db.url")
	if dbURL == "" {
		logger.Error("influxdb url not set")
		return nil
	}

	// Create an influxdb client
	client := influxdb2.NewClient(dbURL, dbToken)

	// validate client connection health
	_, err := client.Health(ctx)
	if err != nil {
		logger.Error("influxdb health check failed", slog.String("err", err.Error()))
		return nil
	}

	return &InfluxdbGateway{
		Logger: logger,
		Client: &client,
	}
}
