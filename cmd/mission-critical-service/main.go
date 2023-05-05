package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BobyMCbobs/sample-ko-monorepo/pkg/common"
)

func httpMustWriteResponse(i int, err error) {
	if err != nil {
		log.Println("error writing response:", err)
	}
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	httpMustWriteResponse(w.Write([]byte("Page not found")))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	httpMustWriteResponse(w.Write([]byte("Healthy")))
}

func getAPINumber(w http.ResponseWriter, r *http.Request) {
	p, _ := rand.Prime(rand.Reader, 64)
	w.WriteHeader(http.StatusOK)
	httpMustWriteResponse(w.Write([]byte(fmt.Sprintf("%v", p))))
}

func getAPIMessage(w http.ResponseWriter, r *http.Request) {
	message := `welcome!`
	w.WriteHeader(http.StatusOK)
	httpMustWriteResponse(w.Write([]byte(message)))
}

type MissionCritialService struct {
	server *http.Server
}

func NewMissionCritialService() *MissionCritialService {
	frontendFolderPath := common.GetServePath()
	appPort := common.GetAppPort()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/_healthz", getHealth)
	mux.HandleFunc("/api/number", getAPINumber)
	mux.HandleFunc("/api/message", getAPIMessage)
	mux.Handle("/", http.FileServer(http.Dir(frontendFolderPath)))
	mux.HandleFunc("/{.*}", pageNotFound)

	handler := common.Logging(mux)
	return &MissionCritialService{
		server: &http.Server{
			Addr:           appPort,
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
