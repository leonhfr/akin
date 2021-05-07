package tags

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type ()

// Tags is an extension for the goldmark
var Tags = &tags{}

type data struct {
	Tags []string
}

var contextKey = parser.NewContextKey()

// Get returns the tag list
func Get(pc parser.Context) []string {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*data)
	return d.Tags
}

type tags struct{}

func (e *tags) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewParser(), 501),
	))
}

type tagsParser struct{}

func NewParser() parser.InlineParser {
	return &tagsParser{}
}

func (s *tagsParser) Trigger() []byte {
	return []byte{'#'}
}

func (s *tagsParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	block.Advance(1)
	line, _ := block.PeekLine()
	if line[0] == byte('#') {
		return nil
	}

	// TODO: stop at the first non alphanumeric character
	tag := strings.Split(string(line), " ")[0]
	ctx := pc.Get(contextKey)
	if ctx != nil {
		d := ctx.(*data)
		d.Tags = append(d.Tags, tag)
		return nil
	}

	d := &data{}
	d.Tags = append(d.Tags, tag)
	pc.Set(contextKey, d)
	return nil
}
