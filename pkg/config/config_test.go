package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tmp_location = "/tmp/config.yaml"
)

func TestLoadAllOptions(t *testing.T) {
	initialize(t)

	content := `backer:
  transfer:
    - source: ./s_test
      destination: ./d_test
    - source: ./1
      destination: ./2
  exclude:
    - log
    - logs
    - "*.log"
    - "node_modules"
  rsync_options:
    - -avAXEWSlHh
    - --no-compress
    - --info=progress2`

	writeToConfig(content)

	actual, err := load(tmp_location)
	assert.Equal(t, nil, err)

	expected := &Config{
		TransferList: []TransferItem{
			{
				Source:      "./s_test",
				Destination: "./d_test",
			},
			{
				Source:      "./1",
				Destination: "./2",
			},
		},
		Exclude:      []string{"log", "logs", "*.log", "node_modules"},
		RsyncOptions: []string{"-avAXEWSlHh", "--no-compress", "--info=progress2"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("expected config does not match the actual")
	}

	clean()
}

func TestLoadEmptyConfig(t *testing.T) {
	initialize(t)

	_, err := load(tmp_location)
	assert.Equal(t, nil, err)

	clean()
}

func TestInvalidConfigFormat(t *testing.T) {
	initialize(t)

	content := `backer:
  exclude:
    - *.log
    - node_modules`

	writeToConfig(content)

	_, err := load(tmp_location)
	assert.Equal(t, ErrInvalidConfigFormat, err)

	clean()
}

func writeToConfig(c string) {
	f, err := os.OpenFile(tmp_location, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = f.WriteString(c)
	if err != nil {
		panic(err)
	}
}

func initialize(t *testing.T) {
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)
}

func clean() {
	err := os.Remove(tmp_location)
	if err != nil {
		panic(err)
	}
	clearAll()
}
