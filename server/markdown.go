package server

import (
	"bytes"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.CJK,
		meta.New(meta.WithTable()),
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.WithClasses(true),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func MarkDownRender(content []byte) template.HTML {
	var buf bytes.Buffer

	if err := md.Convert(content, &buf); err != nil {
		panic(err)
	}
	return template.HTML(buf.String())
}
