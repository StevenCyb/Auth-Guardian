package logging

import (
	"auth-guardian/config"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var lock sync.Mutex

func init() {
	lock.Lock()
	defer lock.Unlock()
	os.Remove("/tmp/now.log")
}

// MapToJSONString build a JSON string out of map
func MapToJSONString(message *map[string]string) string {
	jsonString, _ := json.Marshal(message)
	return string(jsonString)
}

// MapToString build a formatted string out of map
func MapToString(message *map[string]string) string {
	var bString bytes.Buffer
	i := 1
	for k, v := range *message {
		if i == len(*message) {
			bString.WriteString(k + "=" + v)
		} else {
			bString.WriteString(k + "=" + v + ", ")
		}
		i++
	}
	return bString.String()
}

// Log is a helper to reuse code for printing
func Log(message *map[string]string) {
	lock.Lock()
	defer lock.Unlock()

	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			defer f.Close()
			wrt := io.MultiWriter(os.Stdout, f)
			log.SetOutput(wrt)
		}
	}

	output := ""
	if config.LogJSON {
		(*message)["timestamp"] = time.Now().String()
		output = MapToJSONString(message)
	} else {
		output = MapToString(message)
	}
	log.Println(output)
}

// Fatal logs panic
func Fatal(message *map[string]string) {
	(*message)["log_level"] = "panic"
	log.Panic(message)
}

// Error logs an error
func Error(message *map[string]string) {
	if config.LogLevel >= 1 {
		(*message)["log_level"] = "error"
		Log(message)
	}
}

// Warning logs a warning
func Warning(message *map[string]string) {
	if config.LogLevel >= 2 {
		(*message)["log_level"] = "warning"
		Log(message)
	}
}

// Info logs an information
func Info(message *map[string]string) {
	if config.LogLevel >= 3 {
		(*message)["log_level"] = "info"
		Log(message)
	}
}

// Debug logs an debug information
func Debug(message *map[string]string) {
	if config.LogLevel >= 4 {
		(*message)["log_level"] = "debug"
		Log(message)
	}
}
