package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

		w.Write([]byte(fmt.Sprint(entries)))
	} else { // stat is file

	}
}

func (s *Server) Run() error {
	http.HandleFunc("/tree/", s.treeHandler)

	addr := s.Bind + ":" + strconv.Itoa(s.Port)
	fmt.Println("Web Server is available at http://" + addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}
