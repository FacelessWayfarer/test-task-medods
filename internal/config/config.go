package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Conifg struct {
	HTTP struct {
		IP           string        `yaml:"ip" env:"HTTP-IP"`
		Port         int           `yaml:"port" env:"HTTP-PORT"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env:"HTTP-READ-TIMEOUT"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"HTTP-WRITE-TIMEOUT"`
	} `yaml:"http"`
	PostgreSQL struct {
		Username string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"PSQL_PORT" env-required:"true"`
		Database string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
	} `yaml:"postgresql"`
}

const (
	EnvConfigPathName = "CONFIG_PATH"
)

func SetConfig() *Conifg {
	Cfg := &Conifg{}
	var once sync.Once

	once.Do(func() {
		log.Print("initializing config")

		configPath := os.Getenv(EnvConfigPathName)
		log.Println(configPath)
		if configPath == "" {
			log.Fatal("config path is required")
		}

		if err := cleanenv.ReadConfig(configPath, Cfg); err != nil {
			var headerText = "test task"
			errText, _ := cleanenv.GetDescription(Cfg, &headerText)
			log.Print(errText)
			log.Fatal(err)
		}

	})
	return Cfg
}
