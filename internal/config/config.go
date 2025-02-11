package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-default:"./storage/storage.db"`
	HTTPServer  `yaml:"http_server"`
	Admin       `yaml:"admin"`
}

type Admin struct {
	User     string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// Get config path from env or flag
	configPath := fetchConfigPath()
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	// read config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string = os.Getenv("CONFIG_PATH")

	if res == "" {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	return res
}
