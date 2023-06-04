package pkg

import (
	"log"
	"net/http"
)

type LogWriter struct {
	http.ResponseWriter
}

func (w LogWriter) Write(p []byte) {
	_, err := w.ResponseWriter.Write(p)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
