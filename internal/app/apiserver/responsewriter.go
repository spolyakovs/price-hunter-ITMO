package apiserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (writer *responseWriter) WriteHeader(statusCode int) {
	writer.code = statusCode
	writer.ResponseWriter.WriteHeader(statusCode)
}
