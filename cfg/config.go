package cfg

import (
	"fmt"

	"github.com/leonhfr/akin/internal/akin/utils"
)

const version = "0.0.1"

var extensions = []string{".md"}

type CliOptions struct {
	AnkiPort    uint   `long:"anki-port" description:"AnkiConnect port to use" default:"8765" optional:"yes"`
	AnkiAddress string `long:"anki-address" description:"AnkiConnect address to use" default:"localhost" optional:"yes"`
	Path        string `short:"p" long:"path" description:"Path to the local Markdown file or directory" default:"."`
	DryRun      []bool `short:"n" long:"dry-run" description:"Does not apply the changes" optional:"yes"`
	Verbose     []bool `short:"v" long:"verbose" description:"Verbose logs" optional:"yes"`
	Sync        []bool `short:"s" long:"sync" description:"Synchronizes the local Anki collection with AnkiWeb" optional:"yes"`
	Exit        []bool `short:"e" long:"exit" description:"Gracefully exists Anki when done" optional:"yes"`
	Version     []bool `long:"version" description:"Prints current akin version"`
	// cleanUpUsage     = "Deletes from the target all resources not present in the source. Tread carefully."
	// ankiAPIKeyUsage = "Optional AnkiConnect API key."
}

// A Config represents Akin's configuration.
type Config struct {
	CliOptions
	CurrentVersion string
	Extensions     map[string]bool
	URL            string
}

// New creates a new Config.
func New(opts CliOptions) *Config {
	url := fmt.Sprintf("http://%v:%v/", opts.AnkiAddress, opts.AnkiPort)

	return &Config{
		CliOptions:     opts,
		CurrentVersion: version,
		Extensions:     utils.MakeSet(extensions),
		URL:            url,
	}
}
