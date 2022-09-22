package back

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	command        []string
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

	for _, level := range r.transferLevels {
		var cmd []string
		cmd = append(cmd, r.bin)
		cmd = append(cmd, options...)
		cmd = append(cmd, level.source, level.destination)
		level.command = cmd
	}
}

func (r *rsync) run() (string, error) {
	log.Printf("running command: %q", strings.Join(r.command, " "))

	cmd := exec.Command(r.command[0], r.command[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
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
	log.Printf("Syncing from %q to %q", c.Source, c.Destination)
	r := createRsync(c)
	r.build()
	out, err := r.run()
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
