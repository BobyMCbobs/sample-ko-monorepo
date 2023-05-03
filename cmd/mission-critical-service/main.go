package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BobyMCbobs/sample-ko-monorepo/pkg/common"
)

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Page not found")
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Healthy")
}

func getAPINumber(w http.ResponseWriter, r *http.Request) {
	p, _ := rand.Prime(rand.Reader, 64)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("%v", p))
}

type MissionCritialService struct {
	server *http.Server
}

func NewMissionCritialService() *MissionCritialService {
	frontendFolderPath := common.GetServePath()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/_healthz", getHealth)
	mux.HandleFunc("/api/number", getAPINumber)
	mux.Handle("/", http.FileServer(http.Dir(frontendFolderPath)))
	mux.HandleFunc("/{.*}", pageNotFound)

	handler := common.Logging(mux)
	return &MissionCritialService{
		server: &http.Server{
			Addr:           ":8080",
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (w *MissionCritialService) Run() {
	log.Printf("Listening on HTTP port '%v'\n", w.server.Addr)
	log.Fatal(w.server.ListenAndServe())
}

func main() {
	NewMissionCritialService().Run()
}
