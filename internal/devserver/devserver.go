package devserver

import (
	"log"
	"net/http"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/olahol/melody"
)

type devServer struct {
	projectPathAbs string
	watcher        *fsnotify.Watcher
	mux            *http.ServeMux
	m              *melody.Melody
}

func new(opt Options) (*devServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	m := melody.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Println("an application has connected to this dev server!")
	})

	return &devServer{
		projectPathAbs: opt.ProjectPathAbs,
		watcher:        watcher,
		m:              m,
		mux:            mux,
	}, nil
}

func (server devServer) start() error {
	err := server.watcher.Add(path.Join(server.projectPathAbs, "build"))
	if err != nil {
		return err
	}

	go server.listenForFileEvents()

	return http.ListenAndServe(":8759", server.mux)
}

func (server devServer) close() {
	server.close()
}

func (server devServer) listenForFileEvents() {
	var timer *time.Timer

	for {
		select {
		case event, ok := <-server.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.NewTimer(2 * time.Second)
				go func() {
					<-timer.C
					server.onBinaryUpdated()
				}()
			}

		case err, ok := <-server.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (server devServer) onBinaryUpdated() {
	server.m.Broadcast([]byte("hotRestart"))
}
