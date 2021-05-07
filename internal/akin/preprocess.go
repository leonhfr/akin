package akin

import (
	"context"
	"time"

	"github.com/leonhfr/akin/internal/akin/models"
	"github.com/leonhfr/akin/internal/akin/utils"
)

// TODO, Verbose, DryRun
func (akin *Akin) Preprocess() error {
	akin.mutex.Lock()
	defer akin.mutex.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Check version.
	_, err := akin.client.CheckVersion(ctx)
	if err != nil {
		return err
	}

	// Load deck list.
	decks, err := akin.client.Decks(ctx)
	if err != nil {
		return err
	}
	akin.deckMap = utils.MakeSet(decks)

	medias, err := akin.client.GetMedias(ctx, "*")
	if err != nil {
		return err
	}
	akin.mediaMap = utils.MakeSet(medias)

	// Sync models.
	models, err := akin.client.Models(ctx)
	if err != nil {
		return err
	}
	toCreate := missingModels(models)
	err = akin.createModels(toCreate)
	if err != nil {
		return err
	}

	return nil
}

func (akin *Akin) createModels(missing []string) error {
	for _, model := range missing {
		input := models.ModelInput(model)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		err := akin.client.CreateModel(ctx, input)
		cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

func missingModels(current []string) []string {
	var missing []string
	set := utils.MakeSet(current)
	for _, model := range models.MODELS {
		_, ok := set[model]
		if !ok {
			missing = append(missing, model)
		}
	}
	return missing
}
