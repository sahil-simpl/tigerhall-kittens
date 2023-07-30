package app

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"tigerhallKittens/app/lib/db"
	"tigerhallKittens/app/lib/logger"
)

var Env Config

const (
	EnvDevelopment      = "development"
	EnvStaging          = "staging"
	EnvProduction       = "production"
	EnvQualityAssurance = "qa"
)

// Config specifies Basic configuration we need to support the basic functionalities mentioned in the assignment.
type Config struct {
	Environment string `mapstructure:"ENV"`
	Port        string `mapstructure:"PORT"`
	// Database Environment Variable
	DatabaseURL            string `mapstructure:"DATABASE_URL"`
	DatabaseMinConnections string `mapstructure:"DB_MIN_CONNECTIONS"`
	DatabaseMaxConnections string `mapstructure:"DB_MAX_CONNECTIONS"`
	// Service Environment Variable, This will be used to perform auth mentioned in the assignment.
	ServiceID  string `mapstructure:"SERVICE_ID"`
	ServiceKey string `mapstructure:"SERVICE_KEY"`
	// Allowed endpoints to access
	AllowedOrigins string `mapstructure:"ALLOWED_ORIGINS"`
}

// SetupLogger configures the logger in various modes as per the environment.
func SetupLogger(env string) {
	switch env {
	case EnvDevelopment:
		logger.Init(logger.DEBUG)
	case EnvStaging, EnvQualityAssurance:
		logger.Init(logger.INFO)
		fallthrough
	case EnvProduction:
		logger.Init(logger.INFO)
		// TODO: Add Integration to some alerting tool if required by the interviewer.
	}
}

// LoadEnv load the environment variables required for running the service.
func LoadEnv() error {
	// in case of ENV not set or set to development, read development.env
	if os.Getenv("ENV") == EnvDevelopment || os.Getenv("ENV") == "" {
		viper.SetConfigFile("./development.env")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	viper.AutomaticEnv()
	bindEnvs(Env)

	return viper.Unmarshal(&Env)
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

// SetupDBConnection establishes a connection to the DB.
func SetupDBConnection(ctx context.Context) {
	minConnections, err := strconv.Atoi(Env.DatabaseMinConnections)
	if err != nil {
		logger.E(ctx, err, "Invalid DB_MIN_CONNECTIONS", logger.Field("error", err.Error()))
		panic(err)
	}

	maxConnections, err := strconv.Atoi(Env.DatabaseMaxConnections)
	if err != nil {
		logger.E(ctx, err, "Invalid DB_MAX_CONNECTIONS", logger.Field("error", err.Error()))
		panic(err)
	}

	err = db.Connect(Env.DatabaseURL, minConnections, maxConnections)
	if err != nil {
		logger.E(ctx, err, "Failed connecting to database", logger.Field("error", err))
		panic(err)
	}

	logger.I(ctx, "Established connection to database")
}
