package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	log_location = ".config/backer/backer.log"
)

func Log(s ...interface{}) error {
	f, err := os.OpenFile(logLocation(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range s {
		_, err := f.WriteString(fmt.Sprintf("%q\n", v))
		if err != nil {
			return err
		}
	}

	return nil
}

func logLocation() string {
	home, err := homedir.Dir()

	if err != nil {
		panic(err)
	}

	return filepath.Join(home, log_location)
}
