package back

import (
	"reflect"
	"strings"
	"testing"

	"github.com/omarahm3/backer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	r := createRsync(&config.Config{
		TransferList: []config.TransferItem{
			{
				Source:      "./s_test",
				Destination: "./d_test",
			},
			{
				Source:      "./1",
				Destination: "./2",
			},
		},
		Exclude:      []string{"*.log", "node_modules"},
		RsyncOptions: []string{"-avAXEWSlHh"},
	})
	r.build()

	assert.Equal(t, 2, len(r.transferLevels))

	expectedLevels := []transferLevel{
		{
			source:      "./s_test",
			destination: "./d_test",
			command:     []string{"/usr/bin/rsync", "--exclude=*.log", "--exclude=node_modules", "-avAXEWSlHh", "./s_test", "./d_test"},
		},
		{
			source:      "./1",
			destination: "./2",
			command:     []string{"/usr/bin/rsync", "--exclude=*.log", "--exclude=node_modules", "-avAXEWSlHh", "./1", "./2"},
		},
	}

	for i, level := range r.transferLevels {
		if level.source != expectedLevels[i].source {
			t.Errorf("expected source %q does not equal actual source %q on level %d", level.source, expectedLevels[i].source, i)
		}

		if level.destination != expectedLevels[i].destination {
			t.Errorf("expected destination %q does not equal actual destination %q on level %d", level.destination, expectedLevels[i].destination, i)
		}

		if !reflect.DeepEqual(level.command, expectedLevels[i].command) {
			t.Errorf("expected command %q does not equal actual command %q on level %d", strings.Join(level.command, ", "), strings.Join(expectedLevels[i].command, ", "), i)
		}
	}
}
