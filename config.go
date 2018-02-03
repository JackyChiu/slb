package slb

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	Port  int      `json:"port"`
	Hosts []string `json:"hosts"`
}

func ParseConfig(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Errorf("given path doesn't exist: %v", configPath)
	}

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return config, errors.Wrapf(err, "unable to read json file")
}

func MustParseConfig(configPath string) Config {
	config, err := ParseConfig(configPath)
	if err != nil {
		panic(err)
	}
	return config
}
