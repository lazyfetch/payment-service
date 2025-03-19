package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string         `yaml:"env" env-default:"dev"`
	GRPC      GRPCConfig     `yaml:"grpc"`
	Postgres  PostgresConfig `yaml:"postgres"`
	Webhook   WebhookConfig  `yaml:"webhook"`
	RoboKassa Robokassa      `yaml:"robokassa"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	DBname   string `yaml:"dbname"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type WebhookConfig struct {
	Port int `yaml:"port" env-default:"8080"`
}

type Robokassa struct {
	MerchantLogin string `yaml:"merchantlogin"`
	Password      string `yaml:"password1"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
