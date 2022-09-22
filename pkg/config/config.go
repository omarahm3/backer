package config

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/mitchellh/go-homedir"
)

var (
	ErrInvalidConfigFormat  = errors.New("invalid config format")
	ErrCouldNotDecodeConfig = errors.New("could not decode config file")
	ErrCouldNotCreateConfig = errors.New("could not create config file")
)

type Config struct {
	Source      string `mapstructure:"source"`
	Destination string `mapstructure:"destination"`

	// TODO rsync options need to have pre-options and post-options
	RsyncOptions []string `mapstructure:"rsync_options"`
	Exclude      []string `mapstructure:"exclude"`
}

const (
	config_location = ".config/backer/config.yaml"
	config_key      = "backer"
)

func Load() (*Config, error) {
	err := createConfigFile(location())
	if err != nil {
		log.Printf("load:: error creating config file: %s", err.Error())
		return nil, ErrCouldNotCreateConfig
	}
	return load(location())
}

func load(l string) (*Config, error) {
	config.WithOptions(config.Readonly, config.EnableCache)
	config.AddDriver(yamlv3.Driver)
	err := config.LoadFiles(l)
	if err != nil {
		log.Printf("config:: error loading config file: %s", err.Error())
		return nil, ErrInvalidConfigFormat
	}

	var c *Config
	err = config.BindStruct(config_key, &c)
	if err != nil {
		log.Printf("config:: error decoding config file: %s", err.Error())
		return nil, ErrCouldNotDecodeConfig
	}

	return c, nil
}

func createConfigFile(l string) error {
	err := os.MkdirAll(path.Dir(l), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Stat(l)
	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(l)
		return err
	}

	return err
}

func location() string {
	home, err := homedir.Dir()

	if err != nil {
		panic(err)
	}

	return filepath.Join(home, config_location)
}

func clearAll() {
	config.ClearAll()
}
