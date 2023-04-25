package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitLogger() {
	viper.SetDefault("APP_LOG_LEVEL", zerolog.InfoLevel)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	level, err := zerolog.ParseLevel(viper.GetString("APP_LOG_LEVEL"))
	if err != nil {
		log.Fatal().Msgf("Unable to parse level of zerolog: %v", err)
	}

	zerolog.SetGlobalLevel(level)
}
