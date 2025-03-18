package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "config/config.yaml"

type Config struct {
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
}

type Postgres struct {
	DBName  string `yaml:"db_name"`
	DBUser  string `yaml:"db_yser"`
	DBPass  string `yaml:"db_pass"`
	DBHost  string `yaml:"db_host"`
	DBPort  string `yaml:"db_port"`
	SSLMode string `yaml:"db_ssl_mode"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file doesn`t exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can`t read config: %s", err)
	}

	return &cfg
}
