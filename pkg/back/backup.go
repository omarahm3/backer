package back

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/omarahm3/backer/pkg/config"
)

type rsync struct {
	conf        *config.Config
	bin         string
	command     []string
	options     []string
	source      string
	destination string
	excludeList []string
}

func (r *rsync) build() {
	// build source & destination
	r.source = r.conf.Source
	r.destination = r.conf.Destination

	// build exclude list
	excludeList := r.conf.Exclude
	for i, v := range excludeList {
		excludeList[i] = fmt.Sprintf("--exclude=%s", v)
	}
	r.excludeList = excludeList

	// build options
	options := r.conf.RsyncOptions
	options = append(options, r.source, r.destination)

	// build final command
	r.command = []string{r.bin}
	r.command = append(r.command, excludeList...)
	r.command = append(r.command, options...)
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
	return &rsync{
		conf: c,
		bin:  "/usr/bin/rsync",
	}
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
