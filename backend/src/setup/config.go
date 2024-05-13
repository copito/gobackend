package setup

import (
	"log/slog"

	viper "github.com/spf13/viper"
)

func SetupConfig(logger *slog.Logger) {
	viper.SetConfigName("base")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	// Setting prefix for all env variables: SCAFFOLD_ID => viper.Get("ID")
	viper.SetEnvPrefix("DATA_QUALITY")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Failed to load configurations")
		panic(err)
	}

	logger.Info("Configuration loaded successfully...")
}
