package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// An interface for reading game configurations
type IConfig interface {
	// Function that reads a key from the config
	ReadKeyString(string) (string, error)
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

func (s *viperConfig) ReadKeyString(key string) (string, error) {
	val, err := s.ReadKey(key)
	return val.(string), err
}

func (s *viperConfig) WriteKeyValue(key string, val interface{}) error {
	if s.viper_instance == nil {
		return s.noConfigError("no config file loaded")
	}
	s.viper_instance.Set(key, val)
	return nil
}

type mockConfig struct {
}

// Implementation of ReadConfig from IConfig. Reads a yaml configuration file given a path and filename
func (s *mockConfig) ReadConfig(path, file string) error {
	return nil
}

func (s *mockConfig) ReadKey(key string) (interface{}, error) {
	switch key {
	case "redis_host":
		return "redis", nil
	case "redis_pass":
		return "redispassword", nil
	case "redis_port":
		return 6379, nil
	case "host":
		return "mongo", nil
	case "db":
		return "dungeondb", nil
	case "root_user":
		return "dungeonmaster", nil
	case "root_password":
		return "m123123123", nil
	default:
		return "", errors.New("invalid key")
	}
}

func (s *mockConfig) ReadKeyString(key string) (string, error) {
	val, err := s.ReadKey(key)
	return val.(string), err
}

func (s *mockConfig) WriteKeyValue(key string, val interface{}) error {
	return nil
}

func NewIConfig(isMocked bool) IConfig {
	if isMocked {
		return new(mockConfig)
	} else {
		return new(viperConfig)
	}
}
