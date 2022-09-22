package back

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/omarahm3/backer/pkg/config"
)

type rsync struct {
	cmd         string
	options     []string
	source      string
	destination string
	excludeList []string
}

func (r *rsync) run() (string, error) {
	args := r.options
	args = append(args, r.source, r.destination)

	list := []string{r.cmd}
	list = append(list, args...)
	log.Printf("running command: %q", strings.Join(list, " "))

	cmd := exec.Command(r.cmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func createRsync(c *config.Config) rsync {
	return rsync{
		cmd:         "/usr/bin/rsync",
		options:     c.RsyncOptions,
		source:      c.Source,
		destination: c.Destination,
		excludeList: c.Exclude,
	}
}

func Sync(c *config.Config) error {
	log.Printf("Syncing from %q to %q", c.Source, c.Destination)
	r := createRsync(c)
	out, err := r.run()
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
