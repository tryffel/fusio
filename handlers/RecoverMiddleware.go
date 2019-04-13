package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware Recovers if handler panics
func (h *Handler) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		automatic := AutomaticWriter{
			ResponseWriter: w,
		}

		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Panic in rest handler: ", r.RequestURI, ": ", err)
				debug.PrintStack()
				h.Metrics.CounterIncrease("http_panic_recovered", 1)
				JsonErrorResponse(w, "Internal error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(&automatic, r)

		// Make sure each request gets some response. If no response, send internal error
		if !automatic.responded {
			logrus.Errorf("Request didn't receive any response at %s %s", r.Method, r.RequestURI)
			JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		}

	})
}

type AutomaticWriter struct {
	http.ResponseWriter
	responded bool
}

func (a *AutomaticWriter) WriteHeader(status int) {
	a.ResponseWriter.WriteHeader(status)
	a.responded = true
}

func (a *AutomaticWriter) Write(b []byte) (int, error) {
	a.responded = true
	return a.ResponseWriter.Write(b)
}
