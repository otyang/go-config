package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/alexflint/go-arg"
	fromfile "github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

// LoadFromFile reads configuration data from one or more files and merges
// them into a single configuration object. It utilizes struct tags to interpret
// the data.
//
// Supported file extensions are: YAML, JSON, TOML, ENV, EDN
//
// Example usage:
//
//	var myConfig MyConfig
//	err := LoadFromFile(&myConfig, "config1.json", "config2.yaml")
//	if err != nil {
//	   // handle error
//	}
func LoadFromFile[T any](config *T, filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("config error: file path not provided")
	}

	for _, filePath := range filePaths {
		err := fromfile.ReadConfig(filePath, config)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	}

	return nil
}

// LoadFromString parses a configuration string into a config target struct of type T.
// It supports the following formats: JSON, TOML, YAML
//
// Example usage:
//
//	var myConfig MyConfig
//	err := LoadFromString(&myConfig, data, "json")
//	if err != nil {
//	   // handle error
//	}
func LoadFromString[T any](config *T, data, format string) error {
	if len(data) == 0 {
		return errors.New("config error: config is empty")
	}

	var err error

	switch format {
	case "json":
		err = json.Unmarshal([]byte(data), &config)
	case "toml":
		_, err = toml.Decode(data, &config)
	case "yaml":
		err = yaml.Unmarshal([]byte(data), &config)
	default:
		return fmt.Errorf("config string parsing error: unknown/unsupported config format '%s'", format)
	}

	if err != nil {
		return fmt.Errorf("config string parsing error: %s", err.Error())
	}

	return nil
}

// LoadFromCLIFlagsOrENV is a convenient way to parse configuration from
// command-line flags or `cli` environment variables.
//
// Example usage:
//
//	var myConfig MyConfig
//	err := LoadFromCLIFlagsOrENV(&myConfig)
//	if err != nil {
//	   // handle error
//	}
func LoadFromCLIFlagsOrENV[T any](config *T) error {
	return arg.Parse(config)
}
