package testservice

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Run the test service
func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "body { background-color: #333333; }")
	})

	mux.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "console.log('I'm a script.');")
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		img, _ := os.Open("testservice/favicon.ico")
		defer img.Close()
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, img)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello from test service")
	})

	mux.HandleFunc("/mirror", func(w http.ResponseWriter, r *http.Request) {
		mirrorData := make(map[string]interface{})

		// Set cookies
		mirrorData["cookies"] = r.Cookies()

		// Set header
		header := make(map[string]interface{})
		for key, valueSlice := range r.Header {
			header[key] = valueSlice
		}
		mirrorData["header"] = header

		// Set body
		body, _ := ioutil.ReadAll(r.Body)
		mirrorData["body"] = body

		// Return mirror data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mirrorData)
	})

	http.ListenAndServe(":3001", mux)
}
