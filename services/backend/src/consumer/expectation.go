package consumer

// func SampleExpectation(ctx context.Context) {
// 	var influxClient *gateway.InfluxdbGateway
// 	influxClient = gateway.NewInfluxdbGateway(ctx)
// 	if influxClient != nil {
// 		logger.Warn("unable to create influxdb client")
// 		return
// 	}

// 	// TODO: Push data to timeseries (for profileEvent)
// 	org := viper.GetString("timeseries_db.org")
// 	bucket := viper.GetString("timeseries_db.bucket")

// 	queryAPI := client.QueryAPI(org)
// 	query := `from(bucket: "db")
//             |> range(start: -10m)
//             |> filter(fn: (r) => r._measurement == "measurement1")`
// 	results, err := queryAPI.Query(context.Background(), query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for results.Next() {
// 		fmt.Println(results.Record())
// 	}
// 	if err := results.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// }
