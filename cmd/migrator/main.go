// Короче тут надо разобраться что с конфигом за приколы, ибо эта херня не хочет нормально работать.

package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	DBname   string `yaml:"dbname"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

func loadConfig(path string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func newMigrator(cfg *Config, migrationsPath string) (*migrate.Migrate, error) {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return nil, err
	}

	return m, nil

}

func main() {
	var migrationsPath, configPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations file")
	flag.StringVar(&configPath, "config-path", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		panic("config-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	cfg := loadConfig(configPath)

	m, err := newMigrator(cfg, migrationsPath)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
