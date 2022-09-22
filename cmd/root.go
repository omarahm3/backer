package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/omarahm3/backer/pkg/back"
	"github.com/omarahm3/backer/pkg/config"
	"github.com/spf13/cobra"
)

const errInvalidNumberOfArguments = "invalid number of arguments, please enter 2 arguments where first one is source and the next is destination"

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
	config.Log(fmt.Sprintf("received args: %q", strings.Join(args, ",")))
	config.Log(fmt.Sprintf("received rsync options: %q", strings.Join(rsyncOptions, ",")))
	c, err := config.Load()
	check(err)
	err = back.Sync(parseCli(args, cmd, c))
	check(err)
}

func parseCli(args []string, cmd *cobra.Command, c *config.Config) *config.Config {
	if len(args) == 1 || len(args) > 2 {
		fatalPrint(errInvalidNumberOfArguments)
	}

	if len(args) == 2 {
		c.ClearTransfers()
		c.AddTransfer(args[0], args[1])
	}

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
