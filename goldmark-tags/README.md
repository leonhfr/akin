# goldmark-tags

goldmark-tag is an extension for the [goldmark](http://github.com/yuin/goldmark) that allows you to extract tags from a markdown document

## Disclaimer

This package will eventually be moved to its ows repository.

## Usage

### Installation

```
go get github.com/leonhfr/goldmark-tags
```

### Markdown syntax

Tags are alphanumeric strings preceded by a `#` character, without any spaces.

Example:

```md
This is an inline #tag.

Tag list:

#tag1 #tag2 #tag3

Code blocks are not parsed for tags:

`This is not a #tag.`

```

### Access the tags

The code

```go
import (
  "bytes"
  "fmt"

  "github.com/leonhfr/goldmark-tags"
  "github.com/yuin/goldmark"
  "github.com/yuin/goldmark/extension"
  "github.com/yuin/goldmark/parser"
)

var source = `
This is an inline #tag.

Tag list:

#tag1 #tag2 #tag3

Code blocks are not parsed for tags:

\`This is not a #tag.\`
`

func main() {
  markdown := goldmark.New(
    goldmark.WithExtensions(tags.Tags),
  )
  var buf bytes.Buffer
  ctx := parser.NewContext()
  if err := markdown.Convert([]byte(source), &buf, parser.WithContext(ctx)); if err != nil {
    panic(err)
  }

  tagList := tags.Get(ctx)
  fmt.Println(tagList)
}
```

will output

```sh
[tag1 tag2 tag3 notag]
```

The package does not modify the goldmark AST in any way, nor renders the tags.

## License

MIT

## Author

[leonhfr](https://github.com/leonhfr)