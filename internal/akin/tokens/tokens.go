package tokens

import "fmt"

type TokenType int

const (
	TOKEN_TITLE TokenType = iota
	TOKEN_MARKDOWN
	TOKEN_KEY_VALUE
	TOKEN_SYMBOL
	TOKEN_YAML
)

type SymbolType int

const (
	SYMBOL_PERCENT SymbolType = iota
	SYMBOL_HORIZONTAL_RULE
)

type Token struct {
	tokenType  TokenType
	symbolType SymbolType
	key        string
	value      string
	level      int
}

func (token Token) String() string {
	return fmt.Sprintf("token tata %v", token.value)
}

func NewTitle(level int, value string) Token {
	return Token{
		tokenType: TOKEN_TITLE,
		value:     value,
		level:     level,
	}
}

func NewMarkdown(value string) Token {
	return Token{
		tokenType: TOKEN_MARKDOWN,
		value:     value,
	}
}

func NewKeyValue(key, value string) Token {
	return Token{
		tokenType: TOKEN_KEY_VALUE,
		key:       key,
		value:     value,
	}
}

func NewSymbol(symbolType SymbolType) Token {
	return Token{
		tokenType:  TOKEN_MARKDOWN,
		symbolType: symbolType,
	}
}

func NewYaml(yaml string) Token {
	return Token{
		tokenType: TOKEN_YAML,
		value:     yaml,
	}
}

func (token Token) IsTitle() bool {
	return token.tokenType == TOKEN_TITLE
}

func (token Token) IsMarkdown() bool {
	return token.tokenType == TOKEN_MARKDOWN
}

func (token Token) IsKeyValue() bool {
	return token.tokenType == TOKEN_KEY_VALUE
}

func (token Token) IsSymbol() bool {
	return token.tokenType == TOKEN_SYMBOL
}

func (token Token) IsYaml() bool {
	return token.tokenType == TOKEN_YAML
}

func (token Token) Key() (string, error) {
	if !token.IsKeyValue() {
		return "", fmt.Errorf("expected type key value, got %v", token)
	}
	return token.key, nil
}

func (token Token) Value() (string, error) {
	if !token.IsTitle() && !token.IsMarkdown() && !token.IsKeyValue() && !token.IsYaml() {
		return "", fmt.Errorf("expected type title, markdown, keyvalue, or yaml, got %v", token)
	}
	return token.value, nil
}

func (token Token) Level() (int, error) {
	if !token.IsTitle() {
		return 0, fmt.Errorf("expected type title, got %v", token)
	}
	return token.level, nil
}

func (token Token) Symbol() (SymbolType, error) {
	if !token.IsSymbol() {
		return 0, fmt.Errorf("expected type symbol, got %v", token)
	}
	return token.symbolType, nil
}
