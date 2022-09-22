package back

import (
	"bufio"
	"fmt"
	"io"
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
	var wg sync.WaitGroup

	for _, level := range r.transferLevels {
		wg.Add(1)
		config.Log(fmt.Sprintf("running command: %q", strings.Join(level.command, " ")))

		cmd := exec.Command(level.command[0], level.command[1:]...)

		go func(level *transferLevel, cmd *exec.Cmd) {
			defer wg.Done()

			sout, _ := cmd.StdoutPipe()
			go readStd(sout, level)

			serr, _ := cmd.StderrPipe()
			go readStd(serr, level)

			err := cmd.Start()
			if err != nil {
				config.Log(fmt.Sprintf("run:: error occurred running command: %q", err.Error()))
				return
			}

			err = cmd.Wait()
			if err != nil {
				config.Log(fmt.Sprintf("run:: error occurred running command: %q", err.Error()))
				return
			}

		}(level, cmd)
	}
	wg.Wait()

	return nil
}

func readStd(r io.ReadCloser, level *transferLevel) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		fmt.Printf("[%s->%s] > %s\n", level.source, level.destination, s.Text())
	}
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
