package config

import (
	"os"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	PostgresUser      string `mapstructure:"POSTGRES_USER"`
	PostgresHost      string `mapstructure:"POSTGRES_HOST"`
	PostgresPort      int    `mapstructure:"POSTGRES_PORT"`
	PostgresPassword  string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDatabase  string `mapstructure:"POSTGRES_DB"`
	Environment       string `mapstructure:"CHALLENGE_ENV"`
	HTTPServerAddress string `mapstructure:"CHALLENGE_HTTP_SERVER_ADDRESS"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
}

// LoadConfig reads configuration from file or environment variables. Defaults to dev.env in the root path.
func LoadConfig(path string) (Config, error) {
	config := Config{}
	viper.AddConfigPath(path)

	configName := "dev"
	// if CHALLENGE_ENV is set, use that
	if env, ok := os.LookupEnv("CHALLENGE_ENV"); ok {
		configName = env
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
