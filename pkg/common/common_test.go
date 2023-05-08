package common

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func TestLogging(t *testing.T) {
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	var buf bytes.Buffer
	log.SetOutput(&buf)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, req)

	t.Log(buf.String())

	//2023/05/08 16:07:44 GET / HTTP/1.1 <nil>
	re, err := regexp.Compile(`^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} (GET|POST|PATCH|DELETE) [/] HTTP/1.1 <nil>`)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !re.Match(buf.Bytes()) {
		t.Fatalf("error: log output (%v) doesn't match regexp", buf.String())
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	type testCase struct {
		Name          string
		DefaultValue  string
		OSSetValue    string
		ExpectedValue string
	}

	for _, tc := range []testCase{
		{
			Name:          "AAAAA",
			DefaultValue:  "aaaaa",
			OSSetValue:    "",
			ExpectedValue: "aaaaa",
		},
		{
			Name:          "AAAAA",
			DefaultValue:  "aaaaa",
			OSSetValue:    "EEE",
			ExpectedValue: "EEE",
		},
		{
			Name:          "MY_VAR",
			DefaultValue:  "",
			OSSetValue:    "",
			ExpectedValue: "",
		},
	} {
		if tc.OSSetValue == "" {
			if err := os.Unsetenv(tc.Name); err != nil {
				t.Fatalf("error resetting env: %v", err)
			}
		} else {
			if err := os.Setenv(tc.Name, tc.OSSetValue); err != nil {
				t.Fatalf("error: failed to set env; %v", err)
			}
		}
		if output := GetEnvOrDefault(tc.Name, tc.DefaultValue); output != tc.ExpectedValue {
			t.Fatalf("error: env (%v) had unexpected value (%v) instead of (%v)", tc.Name, output, tc.ExpectedValue)
		}
	}
}

func TestGetServePath(t *testing.T) {
	err := os.Unsetenv("KO_DATA_PATH")
	if err != nil {
		t.Fatalf("error resetting env: %v", err)
	}
	if servePath := GetServePath(); servePath != "./cmd/mission-critical-service/kodata" {
		t.Fatalf("error: did not get expected serve path, instead got (%v)", servePath)
	}
	err = os.Setenv("KO_DATA_PATH", "/var/run/ko")
	if err != nil {
		t.Fatalf("error resetting env: %v", err)
	}
	if servePath := GetServePath(); servePath != "/var/run/ko" {
		t.Fatalf("error: did not get expected serve path")
	}
}

func TestGetAppPort(t *testing.T) {
	err := os.Unsetenv("APP_PORT")
	if err != nil {
		t.Fatalf("error resetting env: %v", err)
	}
	if appPort := GetAppPort(); appPort != ":8080" {
		t.Fatalf("error: did not get expected app port, instead got (%v)", appPort)
	}
	err = os.Setenv("APP_PORT", ":8123")
	if err != nil {
		t.Fatalf("error resetting env: %v", err)
	}
	if appPort := GetAppPort(); appPort != ":8123" {
		t.Fatalf("error: did not get expected app port")
	}
}
