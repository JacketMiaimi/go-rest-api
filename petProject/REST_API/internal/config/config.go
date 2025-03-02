package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

// Эта структура будет соответствовать файлу yaml

type Config struct {
	// yaml - какое имя будет имя у соответствующиго параметра
	// yaml - то есть это какое имя у yaml
	// env-required- если забыли установить какй-то параметр, то не запустилось
	// env-required:"true" - без этого параметра программа не запустится
	Env         string `yaml:"environment" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"` // библиотека, которая помогает легче работать с временем
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	// Must - когда функция будет паниковать

	// Getenv - получить окружение
	configPath := os.Getenv("CONFIG_PATH") // Считывать файлы с конфига
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Проверяем существует такой файл
	// os.Stat(configPath) пытается получить информацию о файле.
	// os.IsNotExist(err) - если не найден
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Загружаем конфиг с помощью cleanenv.ReadConfig()
	// cleanenv.ReadConfig(configPath, &cfg) читает YAML-файл и заполняет структуру cfg

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
