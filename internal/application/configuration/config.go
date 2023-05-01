package configuration

import (
	"time"

	"github.com/caarlos0/env"
)

// Application contains data related to application configuration parameters.
type Configuration struct {
	DryRun     bool   `env:"DRY_RUN" envDefault:"false"`
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"production"`
	GinMode    string `env:"GIN_MODE" envDefault:"debug"`

	// Database settings
	DatabaseURI         string        `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	DatabaseMaxOpenConn int           `env:"DATABASE_MAX_OPEN_CONN" envDefault:"10"`
	DatabaseMaxOpenTime time.Duration `env:"DATABASE_MAX_OPEN_TIME" envDefault:""`
	DatabaseMaxIdleConn int           `env:"DATABASE_MAX_IDLE_CONN" envDefault:"10"`
	DatabaseMaxIdleTime time.Duration `env:"DATABASE_MAX_IDLE_TIME" envDefault:""`

	// Kafka settings
	KafkaHost      string `env:"KAFKA_HOST" envDefault:"localhost:29092"`
	KafkaTopic     string `env:"KAFKA_TOPIC" envDefault:"queue"`
	KafkaPortition int    `env:"KAFKA_TOPIC" envDefault:"0"`

	// Autorization settings
	ApiSecret         string        `env:"API_SECRET" envDefault:"example"`
	TokenHourLifespan time.Duration `env:"TOKEN_HOUR_LIFESPAN" envDefault:"1h"`
}

func Load() (Configuration, error) {
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
