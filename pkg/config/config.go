package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// An interface for reading game configurations
type IConfig interface {
	// Function that reads a key from the config
	ReadKey(string) (interface{}, error)
	// Function that writes a key to the config
	WriteKeyValue(string, interface{}) error
	// Function that reads a yaml configuration file given a path and filename
	// Returns non nil error if file could not be read
	ReadConfig(string, string) error
}

type viperConfig struct {
	viper_instance *viper.Viper
}

func (s *viperConfig) noConfigError(reason string) error {
	return errors.New(reason)
}

// Implementation of ReadConfig from IConfig. Reads a yaml configuration file given a path and filename
func (s *viperConfig) ReadConfig(path, file string) error {
	s.viper_instance = viper.New()
	s.viper_instance.SetConfigType("yaml")
	s.viper_instance.AddConfigPath(path)
	s.viper_instance.SetConfigName(file)

	// read in config file and check your errors
	if err := s.viper_instance.ReadInConfig(); err != nil {
		fmt.Println(s.viper_instance.ConfigFileUsed())
		return err
	}

	return nil
}

func (s *viperConfig) ReadKey(key string) (interface{}, error) {
	if s.viper_instance == nil {
		return nil, s.noConfigError("no config file loaded")
	}
	return s.viper_instance.Get(key), nil
}

func (s *viperConfig) WriteKeyValue(key string, val interface{}) error {
	if s.viper_instance == nil {
		return s.noConfigError("no config file loaded")
	}
	s.viper_instance.Set(key, val)
	return nil
}

func NewIConfig() IConfig {
	return new(viperConfig)
}
