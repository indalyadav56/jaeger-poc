package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App  AppConfig  `mapstructure:"app"`
	HTTP HTTPConfig `mapstructure:"http"`
	GRPC GRPCConfig `mapstructure:"grpc"`
	DB   Database   `mapstructure:"database"`
}

type TraceConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type AppConfig struct {
	Name     string      `mapstructure:"name"`
	LogLevel string      `mapstructure:"log_level"`
	Trace    TraceConfig `mapstructure:"trace"`
}

type HTTPConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Database struct {
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
}

type MongoDBConfig struct {
	URI         string `mapstructure:"uri"`
	Database    string `mapstructure:"database"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	AuthSource  string `mapstructure:"auth_source"`
	Timeout     string `mapstructure:"timeout"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
}

func Load() (*Config, error) {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// allow overriding via env vars like APP_NAME, HTTP_PORT
	viper.AutomaticEnv()
	// viper.SetEnvPrefix("PRODUCT_SERVICE")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
