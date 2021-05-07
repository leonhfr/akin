package models

import (
	_ "embed"

	"github.com/leonhfr/anki-connect-go"
)

const (
	MODEL_BASIC = "akin-basic"
	MODEL_DUAL  = "akin-dual"
)

//go:embed css/theme.css
var theme string

var (
	MODELS = []string{MODEL_BASIC, MODEL_DUAL}

	cardTemplate1 = anki.CardTemplateInput{
		Name:  "Card Template 1",
		Front: `{{Front}}`,
		Back: `{{Front}}
		<hr />
		{{Back}}`,
	}

	cardTemplate2 = anki.CardTemplateInput{
		Name:  "Card Template 1",
		Front: `{{Back}}`,
		Back: `{{Back}}
		<hr />
		{{Front}}`,
	}

	basicInput = anki.ModelInput{
		Model:         MODEL_BASIC,
		InOrderFields: []string{"Front", "Back"},
		CSS:           theme,
		CardTemplates: []anki.CardTemplateInput{cardTemplate1},
	}

	dualInput = anki.ModelInput{
		Model:         MODEL_DUAL,
		InOrderFields: []string{"Front", "Back"},
		CSS:           theme,
		CardTemplates: []anki.CardTemplateInput{cardTemplate1, cardTemplate2},
	}
)

func ModelInput(model string) anki.ModelInput {
	switch model {
	case MODEL_BASIC:
		return basicInput
	case MODEL_DUAL:
		return dualInput
	default:
		return basicInput
	}
}
