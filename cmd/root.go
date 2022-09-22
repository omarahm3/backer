package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/omarahm3/backer/pkg/back"
	"github.com/omarahm3/backer/pkg/config"
	"github.com/spf13/cobra"
)

var (
	rsyncOptions []string
	rootCmd      = &cobra.Command{
		Use:   "backer",
		Short: "backup a whole directory based on a set of rules",
		Run:   run,
	}
)

func Init() {
	rootCmd.PersistentFlags().StringSliceVarP(&rsyncOptions, "options", "o", nil, "override rsync options")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("received args: %q", strings.Join(args, ","))
	log.Printf("received rsync options: %q", strings.Join(rsyncOptions, ","))
	c, err := config.Load()
	check(err)
	err = back.Sync(parseCli(args, cmd, c))
	check(err)
}

func parseCli(args []string, cmd *cobra.Command, c *config.Config) *config.Config {
	source, destination := c.Source, c.Destination

	if len(args) > 0 {
		source = args[0]
		args = args[1:]
	}

	if len(args) > 0 {
		destination = args[0]
		args = args[1:]
	}

	c.Source, c.Destination = source, destination

	if rsyncOptions != nil {
		c.RsyncOptions = rsyncOptions
	}

	return c
}

func check(err error) {
	if err == nil {
		return
	}

	fmt.Printf("error occurred: %s\n", err.Error())
	os.Exit(1)
}

func fatalPrint(s string) {
	fmt.Println(s)
	os.Exit(1)
}
