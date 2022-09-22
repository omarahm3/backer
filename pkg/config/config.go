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
	ErrConfigFileIsEmpty    = errors.New("config file is empty")
)

type TransferItem struct {
	Source      string `mapstructure:"source"`
	Destination string `mapstructure:"destination"`
}

type Config struct {
	Source      string `mapstructure:"source"`
	Destination string `mapstructure:"destination"`

	TransferList []TransferItem `mapstructure:"transfer"`
	RsyncOptions []string       `mapstructure:"rsync_options"`
	Exclude      []string       `mapstructure:"exclude"`
}

func (c *Config) ClearTransfers() {
	c.TransferList = []TransferItem{}
}

func (c *Config) AddTransfer(source, destination string) {
	c.TransferList = append(c.TransferList, TransferItem{
		Source:      source,
		Destination: destination,
	})
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
	config.WithOptions(config.EnableCache)
	config.AddDriver(yamlv3.Driver)
	err := config.LoadFiles(l)
	if err != nil {
		log.Printf("config:: error loading config file: %s", err.Error())
		return nil, ErrInvalidConfigFormat
	}

	exists := config.Get(config_key)
	if exists == nil {
		config.Set(config_key, struct{}{})
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
