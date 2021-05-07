package akin

import (
	"context"
	"time"
)

// TODO, Verbose, DryRun
func (akin *Akin) Postprocess() error {
	// if akin.config.Verbose {
	// 	fmt.Println("Syncing Anki...")
	// 	fmt.Println("Exiting Anki...")
	// }
	// if akin.config.DryRun {
	// 	return nil
	// }

	akin.mutex.Lock()
	defer akin.mutex.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := akin.client.Sync(ctx)
	if err != nil {
		return err
	}

	return nil

	// err := akin.client.Exit(ctx)
	// return err
}

// TODO: in dry rune, end with:
// NOTE: Run with "dry run" no changes were made.
