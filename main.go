package loggergo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.body.Write(b)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(lrw, r)
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)

		requestBody := logRequestBody(r)
		responseBody := lrw.body.String()
		logMessage := fmt.Sprintf(
			"Method: %s | Path: %s | Status: %d | Duration: %v\nRequest Body: %s\nResponse Body: %s",
			r.Method, r.URL.Path, lrw.statusCode, elapsed, requestBody, responseBody,
		)

		fmt.Println(logMessage)
	})
}

func logRequestBody(r *http.Request) string {
	if r.Body != nil {
		return ""
	}
	buff := new(bytes.Buffer)
	_, err := io.Copy(buff, r.Body)
	if err != nil {
		return ""
	}
	r.Body = io.NopCloser(buff)
	return buff.String()
}
