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

type config struct {
	postgres postgresConfig `yaml:"postgres"`
}

type postgresConfig struct {
	host     string `yaml:"host" env-default:"localhost"`
	port     int    `yaml:"port"`
	user     string `yaml:"user"`
	dBname   string `yaml:"dbname"`
	password string `yaml:"password"`
}

func loadConfig(path string) *config {
	var cfg config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	fmt.Printf("postgres host: %s\n", cfg.postgres.host)
	fmt.Printf("postgres port: %d\n", cfg.postgres.port)
	fmt.Printf("postgres user: %s\n", cfg.postgres.user)
	fmt.Printf("postgres dbname: %s\n", cfg.postgres.dBname)
	fmt.Printf("postgres password: %s\n", cfg.postgres.password)
	return &cfg
}

func newMigrator(cfg *config, migrationsPath string) (*migrate.Migrate, error) {

	dbURL := "postgres://admin:admin123@localhost:5432/mydatabase?sslmode=disable"

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
