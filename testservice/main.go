package testservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Run the test service
func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(mirrorData)
	})

	http.ListenAndServe(":3001", mux)
}
