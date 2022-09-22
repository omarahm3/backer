package back

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func (r *rsync) build() error {
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
		sPath, err := processSourcePath(level.source)
		if err != nil {
			return err
		}

		level.source = sPath

		dPath, err := processDestinationPath(level.destination)
		if err != nil {
			return err
		}

		level.destination = dPath

		var cmd []string
		cmd = append(cmd, r.bin)
		cmd = append(cmd, options...)
		cmd = append(cmd, level.source, level.destination)
		level.command = cmd
	}

	return nil
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

func processDestinationPath(p string) (string, error) {
	last := p[len(p)-1]
	d := p

	if last == '/' {
		d = filepath.Dir(d)
	}

	d = filepath.Dir(d)
	if !filepath.IsAbs(d) || !isValidPath(d) {
		return "", fmt.Errorf("path %q is not valid", d)
	}

	return p, nil
}

func processSourcePath(p string) (string, error) {
	if !filepath.IsAbs(p) || !isValidPath(p) {
		return "", fmt.Errorf("path %q is not valid", p)
	}

	return p, nil
}

func isValidPath(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func readStd(r io.ReadCloser, level *transferLevel) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		t := fmt.Sprintf("[%s->%s] > %s\n", level.source, level.destination, s.Text())
		fmt.Print(t)
		config.Log(t)
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
	err := r.build()
	if err != nil {
		return err
	}

	err = r.run()
	if err != nil {
		return err
	}

	return nil
}
