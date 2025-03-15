package server

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
		meta.Meta,
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

/*

	if err := md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return err
	}

	metaData := meta.Get(ctx)

	title, ok := metaData["title"]
	if !ok {
		return errors.New("title is not available")
	}

	strTitle, ok := title.(string)
	if !ok {
		return errors.New("title can not assert as string type")
	}

	date, ok := metaData["date"]
	if !ok {
		return errors.New("date is not available")
	}

	strDate, ok := date.(string)
	if !ok {
		return errors.New("date can not assert as string type")
	}

	d, err := time.Parse(time.RFC3339, strDate)
	if err != nil {
		return err
	}
*/

func (a *ArticleTmplContext) MarkDown(content []byte) error {
	var buf bytes.Buffer
	ctx := parser.NewContext()

	if err := md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return err
	}

	metaData := meta.Get(ctx)

	var (
		title sql.NullString
		date  sql.NullTime
	)

	titleStr, ok := metaData["title"].(string)
	title.Valid = ok

	if title.Valid {
		title.String = titleStr
	}

	dateStr, ok := metaData["date"].(string)
	date.Valid = ok

	if date.Valid {
		parsedTime, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return err
		}
		date.Time = parsedTime
	}

	a.Title = title
	a.Date = date
	a.HTML = template.HTML(buf.String())

	return nil
}

// imgRepathize finds all img tag and modifies a path
func imgRepathize(addr, absPath string, content []byte) []byte {
	re := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)

	result := re.ReplaceAllFunc(content, func(m []byte) []byte {
		submatches := re.FindSubmatch(m)
		if len(submatches) < 3 {
			return m
		}

		originalAlt, originalURL := submatches[1], submatches[2]

		if strings.Contains(string(originalURL), "http") {
			return []byte(fmt.Sprintf("![%s](%s)", originalAlt, originalURL))
		}

		dirAbsPath := filepath.Dir(absPath)
		src := filepath.Join(dirAbsPath, string(originalURL))

		newURL := fmt.Sprintf("http://%s/files?src=%s", addr, url.QueryEscape(src))

		return []byte(fmt.Sprintf("![%s](%s)", originalAlt, newURL))
	})

	return result
}
