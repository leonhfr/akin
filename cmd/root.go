package cmd

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/leonhfr/akin/cfg"
	"github.com/leonhfr/akin/internal/akin"
)

var config *cfg.Config

// Execute executes the root command.
func Execute() error {
	a := akin.New(config)
	defer a.Dispose()
	err := a.Preprocess()
	if err != nil {
		return err
	}
	err = a.Process()
	if err != nil {
		return err
	}
	err = a.Postprocess()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	var options cfg.CliOptions
	parser := flags.NewParser(&options, flags.Default)
	parser.Usage = "synchronizes a local repository in Markdown with an Anki collection\n\n  akin [OPTIONS]"
	if _, err := parser.Parse(); err != nil {
		if help() {
			os.Exit(0)
		}
		os.Exit(1)
	}

	config = cfg.New(options)

	if len(config.Version) > 0 {
		version()
	}
}

func help() bool {
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func version() {
	fmt.Printf("current Akin version: %s\n", config.CurrentVersion)
	os.Exit(0)
}
