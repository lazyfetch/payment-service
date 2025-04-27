package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"dev"`
	GRPC        GRPCConfig     `yaml:"grpc"`
	Postgres    PostgresConfig `yaml:"postgres"`
	Redis       RedisConfig    `yaml:"redis"`
	Webhook     WebhookConfig  `yaml:"webhook"`
	Internal    Internal       `yaml:"internal"`
	RateLimiter RateLimiter    `yaml:"rate_limiter"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
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

type Internal struct {
	UserTTL          time.Duration `yaml:"user_TTL"`
	EventSenderTTL   time.Duration `yaml:"event_sender_TTL"`
	MaxNameLength    int           `yaml:"max_name_length"`
	MaxAmount        int64         `yaml:"max_amount"`
	MaxMessageLenght int           `yaml:"max_message_lenght"`
	PaymentService   string        `yaml:"payment_service"`
}

type RateLimiter struct {
	Enabled      bool            `yaml:"enabled"`
	Window       time.Duration   `yaml:"window"`
	MaxRequests  int             `yaml:"max_requests"`
	BanDurations []time.Duration `yaml:"ban_durations"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
