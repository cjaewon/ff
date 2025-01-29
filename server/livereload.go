package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var watcher *fsnotify.Watcher
var connsByPath = make(map[string]map[*websocket.Conn]bool)

// todo: only one goroutine per conn

func (s *Server) liveReloadHandler(w http.ResponseWriter, r *http.Request) {
	// relativePath means "addr/tree/" + relativePath
	relativePath := r.URL.Query().Get("relative-path")
	absPath := filepath.Join(s.RootDirPath, relativePath)
	fmt.Println(absPath)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if connsByPath[absPath] == nil {
		connsByPath[absPath] = map[*websocket.Conn]bool{}
	}

	connsByPath[absPath][conn] = true

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			delete(connsByPath[absPath], conn)

			if len(connsByPath[absPath]) == 0 {
				delete(connsByPath, absPath)
			}

			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			delete(connsByPath[absPath], conn)

			if len(connsByPath[absPath]) == 0 {
				delete(connsByPath, absPath)
			}

			return
		}
	}
}

func watch(dirPath string) error {
	fmt.Println(dirPath)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)

				// todo : fsnotify.Rename must add watch again
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
					fmt.Println(connsByPath)
					// todo: filepath.Dir
					// not only file but also dir
					for conn := range connsByPath[event.Name] {
						fmt.Println("send!")
						fmt.Println(conn.WriteMessage(websocket.TextMessage, []byte("write"+event.Name)))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	count := 0

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if err := watcher.Add(path); err != nil {
				return err
			}
		}

		count += 1

		if count > 30 {
			// todo: panic -> warn -> y/n
			panic("there is too many directory for watching")
		}

		return nil
	})

	return err
}
