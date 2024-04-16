package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	CONFIG_FILE_PATH = fmt.Sprintf("%s/.config/grafana-dashboard-cli/config.yaml", getHomeDir())

	config *Config
)

type TokenMap map[string]string
type Config struct {
	CloudPortal TokenMap `yaml:"cloud-portal"`
	Grafana     TokenMap `yaml:"grafana"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	var err error

	file, err := os.Open(CONFIG_FILE_PATH)
	defer file.Close()

	if err != nil {
		// err = yaml.Unmarshal([]byte(DEFAULT_CONFIG), &cfg)
		return nil, err
	} else {
		err = yaml.NewDecoder(file).Decode(&cfg)
	}

	return &cfg, err
}

func GetConfig() *Config {
	if cfg, err := NewConfig(); err == nil {
		config = cfg
	}
	return config
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}
