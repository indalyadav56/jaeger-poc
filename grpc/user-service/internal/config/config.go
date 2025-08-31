package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	HTTP    HTTPConfig    `mapstructure:"http"`
	GRPC    GRPCConfig    `mapstructure:"grpc"`
	DB      Database      `mapstructure:"database"`
	Clients ClientsConfig `mapstructure:"clients"`
}

type AppConfig struct {
	Name     string      `mapstructure:"name"`
	LogLevel string      `mapstructure:"log_level"`
	Trace    TraceConfig `mapstructure:"trace"`
}

type TraceConfig struct {
	Endpoint string `mapstructure:"endpoint"`
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
	Postgres PostgresDBConfig `mapstructure:"postgres"`
}

type PostgresDBConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  bool   `mapstructure:"ssl_mode"`
}

type ClientsConfig struct {
	ProductService ProductServiceConfig `mapstructure:"product_service"`
}

type ProductServiceConfig struct {
	Target string `mapstructure:"target"`
	Port   int    `mapstructure:"port"`
}

func Load() (*Config, error) {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// allow overriding via env vars like APP_NAME, HTTP_PORT
	viper.AutomaticEnv()

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

func GetDatabaseDSN(cfg *Config) string {
	sslMode := "disable"

	if cfg.DB.Postgres.SSLMode {
		sslMode = "require"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Postgres.Host,
		cfg.DB.Postgres.Port,
		cfg.DB.Postgres.Username,
		cfg.DB.Postgres.Password,
		cfg.DB.Postgres.Name,
		sslMode,
	)
}
