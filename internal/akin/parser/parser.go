package parser

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tags "github.com/leonhfr/akin/goldmark-tags"
	"github.com/leonhfr/akin/internal/akin/models"
	"github.com/leonhfr/akin/internal/akin/tokens"
	"github.com/leonhfr/akin/internal/akin/utils"
	"github.com/leonhfr/anki-connect-go"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Parser struct {
	markdown goldmark.Markdown
	tokens   []*tokens.Token
	titles   []*tokens.Token
	notes    chan<- *anki.NoteInput
	medias   chan<- *anki.StoreMediaInput
	errors   chan<- error
	Deck     string
	Path     string
	index    int
}

func New(t []*tokens.Token, path string, notes chan<- *anki.NoteInput,
	medias chan<- *anki.StoreMediaInput, errors chan<- error) *Parser {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
			tags.Tags,
			extension.Typographer,
			mathjax.MathJax,
		),
		goldmark.WithParserOptions(
			gparser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	return &Parser{
		markdown: markdown,
		tokens:   t,
		titles:   make([]*tokens.Token, 0),
		notes:    notes,
		medias:   medias,
		errors:   errors,
		Deck:     "",
		Path:     path,
		index:    0,
	}
}

func (parser *Parser) Parse() error {
	if parser.isToken() {
		err := parser.parseDeck()
		if err != nil {
			return err
		}
	}

	for parser.isToken() {
		switch {
		case parser.current().IsMarkdown() && len(parser.titles) == 0:
			return fmt.Errorf("expected a title")
		case parser.current().IsTitle():
			parser.appendTitle()
			parser.eat()
		case parser.current().IsMarkdown() && parser.isNextSymbol():
			md, _ := parser.current().Value()
			front := parser.front()
			parser.eat()
			parser.eat()
			backToken := parser.current()
			if !parser.isToken() || !backToken.IsMarkdown() {
				parser.emitNote(front, md)
			} else {
				front = fmt.Sprintf("%v\n\n%v", front, md)
				back, _ := backToken.Value()
				parser.emitNote(front, back)
			}
		case parser.current().IsMarkdown():
			front := parser.front()
			back, _ := parser.current().Value()
			parser.emitNote(front, back)
			parser.eat()
		case parser.current().IsKeyValue():
			key, _ := parser.current().Key()
			value, _ := parser.current().Value()
			parser.emitNote(key, value)
			parser.eat()
		default:
			panic(fmt.Sprintf("unexpected token %s", parser.current()))
		}
	}

	return nil
}

func (parser *Parser) appendTitle() {
	title := parser.current()
	level, _ := title.Level()
	for len(parser.titles) > 0 {
		currentLevel, _ := parser.titles[len(parser.titles)-1].Level()
		if currentLevel > level {
			break
		}
		parser.popTitle()
	}
	parser.titles = append(parser.titles, title)
}

func (parser *Parser) popTitle() *tokens.Token {
	length := len(parser.titles)
	if length == 0 {
		panic("expected to have a length")
	}
	title := parser.titles[length-1]
	parser.titles = parser.titles[:length-1]
	return title
}

func (parser *Parser) current() *tokens.Token {
	return parser.tokens[parser.index]
}

func (parser *Parser) eat() {
	parser.index++
}

func (parser *Parser) emitNote(f, b string) {
	front, tags1 := parser.html(f)
	back, tags2 := parser.html(b)
	tags := utils.MergeStrArr(tags1, tags2)

	parser.notes <- &anki.NoteInput{
		Deck:  parser.Deck,
		Model: models.MODEL_BASIC,
		Fields: anki.FieldsInput{
			Front: front,
			Back:  back,
		},
		Tags: tags,
	}
}

func (parser *Parser) html(md string) (string, []string) {
	var buf bytes.Buffer
	// tags := make([]string, 0)
	ctx := gparser.NewContext()

	err := parser.markdown.Convert([]byte(md), &buf, gparser.WithContext(ctx))
	if err != nil {
		parser.errors <- err
	}

	t := tags.Get(ctx)
	fmt.Println("tags", t)

	return parser.extractMedias(buf.String()), t
}

func (parser *Parser) extractMedias(html string) string {
	re := regexp.MustCompile(`src="([^"]*?)"`)

	return re.ReplaceAllStringFunc(html, func(uri string) string {
		_, err := url.Parse(uri)
		if err != nil {
			// URL, return as is
			return uri
		}

		path := filepath.Join(parser.Path, uri[5:len(uri)-1])
		hash, err := md5File(path)
		if err != nil {
			parser.errors <- err
			return uri
		}

		extension := filepath.Ext(path)
		filename := fmt.Sprintf("%v%v", hash, extension)
		replace := fmt.Sprintf("src=\"%v\"", filename)

		parser.medias <- &anki.StoreMediaInput{
			Filename: filename,
			Path:     path,
		}

		return replace
	})
}

func md5File(path string) (string, error) {
	var res string
	file, err := os.Open(path)
	if err != nil {
		return res, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return res, err
	}
	hashBytes := hash.Sum(nil)[:16]
	res = hex.EncodeToString(hashBytes)
	return res, nil
}

func (parser *Parser) front() string {
	titles := make([]string, 0)
	for _, title := range parser.titles {
		content, _ := title.Value()
		titles = append(titles, content)
	}
	return strings.Join(titles, "\n\n")
}

func (parser *Parser) isNextSymbol() bool {
	if !parser.isNextToken() {
		return false
	}
	return parser.tokens[parser.index+1].IsSymbol()
}

func (parser *Parser) isNextToken() bool {
	return parser.index+1 < len(parser.tokens)
}

func (parser *Parser) isToken() bool {
	return parser.index < len(parser.tokens)
}

func (parser *Parser) parseDeck() error {
	token := parser.current()
	level, err := token.Level()
	if err != nil || level != 1 {
		return fmt.Errorf("expected title of level 1, got %v", token)
	}
	parser.Deck, _ = token.Value()
	parser.eat()
	return nil
}
