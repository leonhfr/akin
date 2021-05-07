package lexer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/leonhfr/akin/internal/akin/tokens"
	"gopkg.in/yaml.v2"
)

type LexerMode int

const (
	LEXER_CLASSIC LexerMode = iota
	LEXER_KEY_VALUE
)

type Lexer struct {
	file    *os.File
	scanner *bufio.Scanner
	tokens  []*tokens.Token
	errors  chan<- error
	config  *lexConfig
	state   lexFn
	line    int
	next    bool
}

func New(path string, errors chan<- error) (*Lexer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	config := newConfig()

	return &Lexer{
		scanner: scanner,
		tokens:  make([]*tokens.Token, 0),
		errors:  errors,
		config:  config,
		state:   lexStart,
		line:    1,
		next:    true,
	}, nil
}

type lexConfig struct {
	mode LexerMode
}

func newConfig() *lexConfig {
	return &lexConfig{
		mode: LEXER_CLASSIC,
	}
}

func (lexer *Lexer) Lex() ([]*tokens.Token, error) {
	defer lexer.file.Close()

	for lexer.isNext() {
		lexer.state = lexer.state(lexer)
	}

	if err := lexer.scanner.Err(); err != nil {
		return nil, err
	}

	return lexer.tokens, nil
}

func (lexer *Lexer) current() string {
	return lexer.scanner.Text()
}

func (lexer *Lexer) eat() bool {
	lexer.line++
	lexer.next = lexer.scanner.Scan()
	return lexer.next
}

func (lexer *Lexer) emit(token tokens.Token) {
	lexer.tokens = append(lexer.tokens, &token)
}

func (lexer *Lexer) emitError(str string) {
	lexer.errors <- fmt.Errorf("%v on line %v", str, lexer.line)
}

func (lexer *Lexer) isNext() bool {
	return lexer.next
}

type lexerYamlFrontMatter struct {
	Mode string `yaml:"mode"`
}

func (lexer *Lexer) setConfig(yml string) {
	var frontMatter lexerYamlFrontMatter
	yaml.Unmarshal([]byte(yml), &frontMatter)

	switch frontMatter.Mode {
	case "classic":
		lexer.config.mode = LEXER_CLASSIC
	case "keyvalue", "kv":
		lexer.config.mode = LEXER_KEY_VALUE
	default:
		lexer.config.mode = LEXER_CLASSIC
	}
}
