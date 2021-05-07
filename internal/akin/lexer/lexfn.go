package lexer

import (
	"regexp"
	"strings"

	"github.com/leonhfr/akin/internal/akin/tokens"
)

type lexFn func(*Lexer) lexFn

var titleRegex = regexp.MustCompile(`^(#+) (.+)$`)

func lexStart(lexer *Lexer) lexFn {
	lexer.eat()
	line := lexer.current()
	if isFrontMatterDelimiter(line) {
		return lexFrontMatter
	}
	return lexContent
}

func lexFrontMatter(lexer *Lexer) lexFn {
	// lexer.emit(tokens.NewSymbol(tokens.SYMBOL_HORIZONTAL_RULE))
	// defer lexer.emit(tokens.NewSymbol(tokens.SYMBOL_HORIZONTAL_RULE))

	yml := ""
	for {
		if !lexer.isNext() {
			// TODO: error
			break
		}
		lexer.eat()

		line := lexer.current()
		if isFrontMatterDelimiter(line) {
			break
		}

		yml += "\n" + line
	}

	lexer.eat()
	// lexer.emit(tokens.NewYaml(yml))
	lexer.setConfig(yml)
	return lexContent
}

func lexContent(lexer *Lexer) lexFn {
	switch lexer.config.mode {
	case LEXER_KEY_VALUE:
		return lexKeyValue
	default:
		return lexClassic
	}
}

func lexClassic(lexer *Lexer) lexFn {
	line := lexer.current()
	switch {
	case isEmptyLine(line):
		return lexEmptyLine
	case isTitle(line):
		return lexTitle
	case isPercent(line):
		return lexPercent
	default:
		return lexMarkdown
	}
}

// TODO: add % support for multiline vocabulary
func lexKeyValue(lexer *Lexer) lexFn {
	line := lexer.current()
	switch {
	case isEmptyLine(line):
		return lexEmptyLine
	case isTitle(line):
		return lexTitle
	default:
		return lexKeyValueLine
	}
}

func lexEmptyLine(lexer *Lexer) lexFn {
	lexer.eat()
	return lexContent
}

func lexMarkdown(lexer *Lexer) lexFn {
	md := lexer.current()
	for {
		lexer.eat()
		if !lexer.isNext() {
			break
		}

		line := lexer.current()
		title := isTitle(line)
		if title {
			break
		}

		symbol := isPercent(line)
		if symbol {
			break
		}

		md += "\n\n" + line
	}
	lexer.emit(tokens.NewMarkdown(md))
	return lexContent
}

func lexKeyValueLine(lexer *Lexer) lexFn {
	line := lexer.current()
	parts := strings.Split(line, "=")

	if len(parts) < 2 {
		lexer.emitError("expected a = symbol")
		return lexContent
	}

	if len(parts) > 2 {
		lexer.emitError("too many = symbols")
		return lexContent
	}

	key, value := parts[0], parts[1]
	lexer.emit(tokens.NewKeyValue(key, value))

	lexer.eat()
	return lexContent
}

func lexPercent(lexer *Lexer) lexFn {
	lexer.emit(tokens.NewSymbol(tokens.SYMBOL_PERCENT))
	lexer.eat()
	return lexContent
}

func lexTitle(lexer *Lexer) lexFn {
	line := lexer.current()
	level, title := parseTitle(line)
	lexer.emit(tokens.NewTitle(level, title))
	lexer.eat()
	return lexContent
}

func isEmptyLine(line string) bool {
	return line == ""
}

func isFrontMatterDelimiter(line string) bool {
	return line == "---"
}

func isPercent(line string) bool {
	return line == "%"
}

func isTitle(line string) bool {
	matches := titleRegex.FindStringSubmatch(line)
	if len(matches) == 0 || len(matches) < 3 {
		return false
	}
	return true
}

func parseTitle(line string) (int, string) {
	matches := titleRegex.FindStringSubmatch(line)
	return len(matches[1]), matches[2]
}
