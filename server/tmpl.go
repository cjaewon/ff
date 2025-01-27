package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

var (
	//go:embed web/template
	tmplFS      embed.FS
	dirTmpl     *template.Template
	articleTmpl *template.Template
)

type DirTmplContext struct {
	RootDirPath string
	RootDirName string
	// Pwdf means print working directory or file.
	Pwdf    []Path
	Entries []File
}

type ArticleTmplContext struct {
	Title      string
	Date       time.Time
	HTML       template.HTML
	IsMarkDown bool
}

func init() {
	b, err := fs.ReadFile(tmplFS, "web/template/dir.html")
	if err != nil {
		panic(err)
	}

	dirTmpl = template.Must(template.New("dir-tmpl").Parse(string(b)))

	b, err = fs.ReadFile(tmplFS, "web/template/article.html")
	if err != nil {
		panic(err)
	}

	articleTmpl = template.Must(template.New("article-tmpl").Parse(string(b)))
}

func (d *DirTmplContext) Write(w http.ResponseWriter) {
	if err := dirTmpl.Execute(w, d); err != nil {
		fmt.Println(err)
	}
}

func (a *ArticleTmplContext) Write(w http.ResponseWriter) {
	if err := articleTmpl.Execute(w, a); err != nil {
		fmt.Println(err)
	}
}

type Path struct {
	Base string
	URL  string
}

type File struct {
	Name      string
	URL       string
	UpdatedAt time.Time
}
