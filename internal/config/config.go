package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/FacelessWayfarer/test-task-medods/pkg/logger"
)

type IConfig interface {
}

type Conifg struct {
	HTTP struct {
		IP           string        `yaml:"ip" env:"HTTP-IP" env-required:"true"`
		Port         int           `yaml:"port" env:"HTTP-PORT" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP-READ-TIMEOUT" env-default:"3s"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP-WRITE-TIMEOUT" env-default:"5s"`
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
	envConfigPathName = "CONFIG_PATH"
)

func SetConfig(logger logger.Logger) (*Conifg, error) {
	Cfg := &Conifg{}

	logger.Println("Initializing config")

	configPath := os.Getenv(envConfigPathName)

	if configPath == "" {
		return nil, fmt.Errorf("config path is required")
	}

	if err := cleanenv.ReadConfig(configPath, Cfg); err != nil {
		var headerText = "test_task"

		errText, _ := cleanenv.GetDescription(Cfg, &headerText)

		return nil, fmt.Errorf("%v : %v", errText, err)
	}

	return Cfg, nil
}
