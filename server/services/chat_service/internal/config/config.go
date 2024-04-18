package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env           string        `yaml:"env" env-default:"local"`
	Storage       StorageConfig `yaml:"storage_cfg" env-required:"true"`
	GRPC          GRPCConfig    `yaml:"grpc" env-required:"true"`
	OtherServices OtherServices `yaml:"other_services" env-required:"true"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type StorageConfig struct {
	RedisURL string `yaml:"redis_url" env-required:"true"`
}

type OtherServices struct {
	AuthService string `yaml:"auth_service_url" env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
