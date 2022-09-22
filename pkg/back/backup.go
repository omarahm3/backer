package back

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/omarahm3/backer/pkg/config"
)

type transferLevel struct {
	command     []string
	source      string
	destination string
}

type rsync struct {
	conf           *config.Config
	bin            string
	excludeList    []string
	transferLevels []*transferLevel
}

func (r *rsync) build() {
	// build exclude list
	excludeList := r.conf.Exclude
	for i, v := range excludeList {
		excludeList[i] = fmt.Sprintf("--exclude=%s", v)
	}
	r.excludeList = excludeList

	// build options
	options := excludeList
	options = append(options, r.conf.RsyncOptions...)

	// build transfer levels
	for _, level := range r.transferLevels {
		var cmd []string
		cmd = append(cmd, r.bin)
		cmd = append(cmd, options...)
		cmd = append(cmd, level.source, level.destination)
		level.command = cmd
	}
}

func (r *rsync) run() error {
	// TODO make sure to run each command on a go routine
	var wg sync.WaitGroup

	for _, level := range r.transferLevels {
		wg.Add(1)
		config.Log(fmt.Sprintf("running command: %q", strings.Join(level.command, " ")))

		go func(level *transferLevel) {
			defer wg.Done()
			cmd := exec.Command(level.command[0], level.command[1:]...)
			out, _ := cmd.CombinedOutput()

			fmt.Println(string(out))
		}(level)

	}
	wg.Wait()

	return nil
}

func createRsync(c *config.Config) *rsync {
	r := &rsync{
		conf:           c,
		bin:            "/usr/bin/rsync",
		transferLevels: []*transferLevel{},
	}

	for _, level := range c.TransferList {
		r.transferLevels = append(r.transferLevels, &transferLevel{
			source:      level.Source,
			destination: level.Destination,
		})
	}

	return r
}

func Sync(c *config.Config) error {
	r := createRsync(c)
	r.build()
	err := r.run()
	if err != nil {
		return err
	}

	return nil
}
