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
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)

	content := `backer:
  source: ./s_test
  destination: ./d_test
  exclude:
    - log
    - logs
    - "*.log"
    - "node_modules"
  rsync_options:
    - -avAXEWSlHh
    - --no-compress
    - --info=progress2`

	fillConfigFile(content)

	actual, err := load(tmp_location)
	assert.Equal(t, nil, err)

	expected := &Config{
		Source:       "./s_test",
		Destination:  "./d_test",
		Exclude:      []string{"log", "logs", "*.log", "node_modules"},
		RsyncOptions: []string{"-avAXEWSlHh", "--no-compress", "--info=progress2"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("expected config does not match the actual")
	}

	clean(tmp_location)
}

func TestLoadEmptyConfig(t *testing.T) {
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)

	_, err = load(tmp_location)
	assert.Equal(t, nil, err)

	clean(tmp_location)
}

func TestInvalidConfigFormat(t *testing.T) {
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)

	content := `backer:
  exclude:
    - *.log
    - node_modules`

	fillConfigFile(content)

	_, err = load(tmp_location)
	assert.Equal(t, ErrInvalidConfigFormat, err)

	clean(tmp_location)
}

func fillConfigFile(c string) {
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

func clean(f string) {
	err := os.Remove(f)
	if err != nil {
		panic(err)
	}
	clearAll()
}
