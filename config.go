package slb

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// Config models the configuations that is avaiable for the balancer.
type Config struct {
	Port  int      `json:"port"`
	Hosts []string `json:"hosts"`
}

// ParseConfig reads a json file and maps it to a Config object.
func ParseConfig(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Errorf("given path doesn't exist: %v", configPath)
	}

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return config, errors.Wrapf(err, "unable to read json file")
}

// MustParseConfig parse and config file and panics on failure.
func MustParseConfig(configPath string) Config {
	config, err := ParseConfig(configPath)
	if err != nil {
		panic(err)
	}
	return config
}
