package config

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DataBaseURL string `toml:"database_url"`
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath,
		"config-path",
		"./configs/apiserver.toml",
		"path for config file")
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
	}
}

func (c *Config) ParseFlags() {
	flag.Parse()

	if _, err := toml.DecodeFile(configPath, c); err != nil {
		log.Fatal(err)
	}
}
