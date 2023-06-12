package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BobyMCbobs/sample-ko-monorepo/pkg/common"
)

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Page not found")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		pageNotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from some kinda web thingy")
}

type WebThingy struct {
	server *http.Server
}

func NewWebThingy() *WebThingy {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	handler := common.Logging(mux)
	return &WebThingy{
		server: &http.Server{
			Addr:              ":8080",
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func (w *WebThingy) Run() {
	log.Printf("Listening on HTTP port '%v'\n", w.server.Addr)
	log.Fatal(w.server.ListenAndServe())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := w.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}

func main() {
	NewWebThingy().Run()
}
