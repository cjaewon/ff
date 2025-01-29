package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
)

type Server struct {
	Port        int
	Watch       bool
	Bind        string
	RootDirName string
	RootDirPath string
}

/*
localhost:1234/tree/
*/

func (s *Server) treeHandler(w http.ResponseWriter, r *http.Request) {
	relativePath := strings.TrimPrefix(r.URL.Path, "/tree/")
	absPath := filepath.Join(s.RootDirPath, relativePath)

	stat, err := os.Stat(absPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if stat.IsDir() {
		entries, err := os.ReadDir(absPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		d := &DirTmplContext{
			RootDirPath: s.RootDirPath,
			RootDirName: s.RootDirName,
			LiveReload:  s.Watch,
		}
		d.Pwdf = append(d.Pwdf, Path{
			Base: s.RootDirName,
			URL:  "/tree/",
		})

		for _, entity := range entries {
			info, err := entity.Info()
			if err != nil {
				fmt.Println(err)
				return
			}

			url, err := url.JoinPath(r.URL.Path, info.Name())
			if err != nil {
				fmt.Println(err)
				return
			}

			d.Entries = append(d.Entries, File{
				Name:      info.Name(),
				URL:       url,
				UpdatedAt: info.ModTime(),
			})
		}

		relativePaths := strings.Split(relativePath, "/")

		for i, base := range relativePaths {
			url, err := url.JoinPath("/tree/", strings.Join(relativePaths[:i], "/"), base)
			if err != nil {
				fmt.Println(err)
				return
			}

			d.Pwdf = append(d.Pwdf, Path{
				Base: base,
				URL:  url,
			})
		}

		d.Write(w)
	} else { // stat is file
		ext := filepath.Ext(absPath)

		if !(ext == ".md" || ext == ".markdown") {
			a := &ArticleTmplContext{
				Title: "지원하지 않는 파일 형태 입니다.",
				Date:  time.Now(),
				HTML: template.HTML(heredoc.Docf(`
					<p><strong>"%s"</strong> 를 읽을 수 없습니다.</p>
				`, filepath.Base(absPath))),
				IsMarkDown: false,
				LiveReload: s.Watch,
			}

			a.Write(w)
			return
		}

		b, err := os.ReadFile(absPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		a := &ArticleTmplContext{
			IsMarkDown: true,
		}

		addr := s.Bind + ":" + strconv.Itoa(s.Port)
		b = imgRepathize(addr, absPath, b)

		if err := a.MarkDown(b); err != nil {
			fmt.Println(err)
			return
		}

		a.Write(w)
	}
}

func (s *Server) filesHandler(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	if src == "" {
		fmt.Println("src is not found")
		return
	}

	http.ServeFile(w, r, src)
}

var (
	//go:embed web/assets
	webAssetsFS embed.FS
	assetsFs, _ = fs.Sub(webAssetsFS, "web")
)

func (s *Server) Run() error {
	http.HandleFunc("/tree/", s.treeHandler)
	http.HandleFunc("/files", s.filesHandler)

	if s.Watch {
		go watch()
		http.HandleFunc("/livereload", s.liveReloadHandler)
	}

	http.Handle("/assets/", http.FileServer(http.FS(assetsFs)))

	addr := s.Bind + ":" + strconv.Itoa(s.Port)
	fmt.Println("Web Server is available at http://" + addr + "/tree")

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}
