package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Conifg struct {
	HTTP struct {
		IP           string        `yaml:"ip" env:"HTTP-IP" env-required:"true"`
		Port         int           `yaml:"port" env:"HTTP-PORT" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"readtimeout" env:"HTTP-READ-TIMEOUT" env-default:"3s"`
		WriteTimeout time.Duration `yaml:"writetimeout" env:"HTTP-WRITE-TIMEOUT" env-default:"5s"`
	} `yaml:"http"`
	PostgreSQL struct {
		Username string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"PSQL_PORT" env-default:"6060"`
		Database string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
	} `yaml:"postgresql"`
}

const (
	EnvConfigPathName = "CONFIG_PATH"
)

func SetConfig() *Conifg {
	Cfg := &Conifg{}

	log.Print("Initializing config")

	configPath := os.Getenv(EnvConfigPathName)

	if configPath == "" {
		log.Fatal("config path is required")
	}

	if err := cleanenv.ReadConfig(configPath, Cfg); err != nil {

		var headerText = "test task"

		errText, _ := cleanenv.GetDescription(Cfg, &headerText)

		log.Print(errText)

		log.Fatal(err)
	}

	return Cfg
}
