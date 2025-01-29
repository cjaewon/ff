package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var watchMap = make(map[*websocket.Conn]watchInfo)

type watchInfo struct {
	Time    time.Time
	AbsPath string
}

// todo: only one goroutine per conn
func (s *Server) liveReloadHandler(w http.ResponseWriter, r *http.Request) {
	// relativePath means "addr/tree/" + relativePath
	relativePath := r.URL.Query().Get("relative-path")
	absPath := filepath.Join(s.RootDirPath, relativePath)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	stat, err := os.Stat(absPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		delete(watchMap, conn)

		return nil
	})

	watchMap[conn] = watchInfo{
		Time:    stat.ModTime(),
		AbsPath: absPath,
	}
	// todo: watchMap 동시성 처리 mutex
}

func watch() error {
	for {
		for conn, info := range watchMap {
			stat, err := os.Stat(info.AbsPath)
			if err == os.ErrNotExist {
				continue
			} else if err != nil {
				return err
			}

			if stat.ModTime().After(info.Time) {
				if _, ok := watchMap[conn]; ok {
					watchMap[conn] = watchInfo{
						Time:    stat.ModTime(),
						AbsPath: info.AbsPath,
					}
				}

				conn.WriteJSON(Message{
					Command: "refresh",
				})
			}
		}

		time.Sleep(time.Second)
	}

}

type Message struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}
