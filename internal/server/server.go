package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Server struct {
	DistDir   string
	BuildFn   func() error
	WatchDirs []string
	Addr      string

	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/_reload", s.handleSSE)
	mux.Handle("/", http.FileServer(http.Dir(s.DistDir)))
	return mux
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 1)
	s.addClient(ch)
	defer s.removeClient(ch)

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) addClient(ch chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clients == nil {
		s.clients = make(map[chan struct{}]struct{})
	}
	s.clients[ch] = struct{}{}
}

func (s *Server) removeClient(ch chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, ch)
}

func (s *Server) notifyClients() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (s *Server) Start() error {
	if s.BuildFn == nil {
		return fmt.Errorf("BuildFn is required")
	}

	// Initial build
	log.Println("Running initial build...")
	if err := s.BuildFn(); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}

	// Start file watcher
	if err := s.startWatcher(); err != nil {
		return fmt.Errorf("start watcher: %w", err)
	}

	addr := s.Addr
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Serving at http://localhost%s\n", addr)
	log.Println("Watching for changes...")
	return http.ListenAndServe(addr, s.Handler())
}

func (s *Server) startWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					log.Printf("Change detected: %s", event.Name)
					if err := s.BuildFn(); err != nil {
						log.Printf("Build error: %v", err)
						continue
					}
					s.notifyClients()
					log.Println("Rebuild complete, reloading browsers...")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	for _, dir := range s.WatchDirs {
		if err := addRecursive(watcher, dir); err != nil {
			log.Printf("Warning: could not watch %s: %v", dir, err)
		}
	}

	return nil
}

func addRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}
