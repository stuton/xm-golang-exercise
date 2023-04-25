package utils

import (
	"flag"
	"log"
	"path"

	"github.com/spf13/viper"
)

const (
	LocalConfigEnabledKey = "local-config-enabled"
	LocalConfigPathKey    = "local-config-path"
	LocalConfigNameKey    = "local-config-name"
)

const (
	LocalConfigNameDefault = "local.env"
)

var (
	LocalConfigEnabled = false
	LocalConfigPath    = "../../config/"
	LocalConfigName    = LocalConfigNameDefault
)

func init() {
	flag.BoolVar(&LocalConfigEnabled, LocalConfigEnabledKey, LocalConfigEnabled, "enable development config ENV parser")
	flag.StringVar(&LocalConfigPath, LocalConfigPathKey, LocalConfigPath, "path to *.env folder")
	flag.StringVar(&LocalConfigName, LocalConfigNameKey, LocalConfigName, "values ENV file to read config from")
}

func InitConfig() {

	flag.Parse()
	viper.AutomaticEnv()

	if LocalConfigEnabled {
		viper.SetConfigType("env")
		viper.SetConfigFile(path.Join(LocalConfigPath, LocalConfigName))
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatalf("Failing to find configuration file: %v", err)
			} else {
				log.Fatalf("Unable to read configuration file: %v", err)
			}
		}
	}
}
