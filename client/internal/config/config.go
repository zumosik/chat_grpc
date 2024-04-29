package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
}

func MustLoadConfig() *Config {
	pathToCfg := mustGetPath()

	var cfg Config

	err := cleanenv.ReadConfig(pathToCfg, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

// mustGetPath is a helper function that returns the path to the configuration file.
// Priority: flag > env > default
func mustGetPath() string {
	var s string
	flag.StringVar(&s, "config", "", "path to the configuration file")
	flag.Parse()

	if s == "" {
		s = os.Getenv("CONFIG")
	}

	if s == "" {
		s = "./configs/local.yml"
	}

	return s
}
