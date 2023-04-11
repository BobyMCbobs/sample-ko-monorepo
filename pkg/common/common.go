package common

import (
	"log"
	"net/http"
	"os"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v", r.Method, r.URL, r.Proto, r.Response)
		next.ServeHTTP(w, r)
	})
}

func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output, ok := os.LookupEnv(envName)
	if !ok {
		output = defaultValue
	}
	return output
}

func GetServePath() string {
	return GetEnvOrDefault("KO_DATA_PATH", "./cmd/mission-critical-service/kodata")
}
