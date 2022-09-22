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

func TestLoadEmptyConfig(t *testing.T) {
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)

	_, err = load(tmp_location)
	assert.Equal(t, nil, err)

	clean(tmp_location)
}

func TestLoadExcludeConfig(t *testing.T) {
	err := createConfigFile(tmp_location)
	assert.Equal(t, nil, err)

	content := `backer:
  exclude:
    - "*.log"
    - node_modules`

	fillConfigFile(content)

	c, err := load(tmp_location)
	actual := c.Exclude
	expected := []string{"*.log", "node_modules"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected list %q, got %q", expected, actual)
	}

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
