package handlers

import (
	"net/http"
	"time"
)

type LoggingWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (l *LoggingWriter) WriteHeader(status int) {
	l.length = 0
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *LoggingWriter) Write(b []byte) (int, error) {
	l.length = len(b)
	if l.status == 0 {
		l.status = 200
	}
	return l.ResponseWriter.Write(b)
}

// LogginMiddlware Provide logging for requests
func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := &LoggingWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(logger, r)

		h.Metrics.CounterIncrease("http_total_requests", 1)

		duration := time.Since(start).String()
		verb := r.Method
		url := r.RequestURI

		fields := make(map[string]interface{})
		fields["verb"] = verb
		fields["request"] = url
		fields["duration"] = duration
		fields["status"] = logger.status
		fields["length"] = logger.length
		h.RequestsLog.WithFields(fields).Infof("")
	})
}
