package main

import (
	"fmt"
	"os"

	"github.com/leonhfr/akin/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if s := err.Error(); s != "" {
			fmt.Printf("akin: %s\n", s)
		}
		os.Exit(1)
	}
}
