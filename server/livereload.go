package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var watchMap = make(map[*websocket.Conn]watchInfo)
var watchMapMutex sync.Mutex

type watchInfo struct {
	Time    time.Time
	AbsPath string
}

func (s *Server) liveReloadHandler(w http.ResponseWriter, r *http.Request) {
	// relativePath means "addr/tree/" + relativePath
	relativePath := r.URL.Query().Get("relative-path")
	absPath := filepath.Join(s.RootDirPath, relativePath)

	stat, err := os.Stat(absPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		watchMapMutex.Lock()
		delete(watchMap, conn)
		watchMapMutex.Unlock()

		return nil
	})

	watchMapMutex.Lock()
	watchMap[conn] = watchInfo{
		Time:    stat.ModTime(),
		AbsPath: absPath,
	}
	watchMapMutex.Unlock()

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}
		}
	}()
}

func watch() {
	for {
		watchMapMutex.Lock()

		for conn, info := range watchMap {
			stat, err := os.Stat(info.AbsPath)
			if errors.Is(err, os.ErrNotExist) {
				// todo: conn WriteJSON becase of showing file is not existed
				// made some 404 page with no websocket connect
				continue
			} else if err != nil {
				fmt.Println(err)
				continue
			}

			if stat.ModTime().After(info.Time) {
				watchMap[conn] = watchInfo{
					Time:    stat.ModTime(),
					AbsPath: info.AbsPath,
				}

				conn.WriteJSON(Message{
					Command: "refresh",
				})
			}
		}
		watchMapMutex.Unlock()

		time.Sleep(time.Second)
	}

}

type Message struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}
