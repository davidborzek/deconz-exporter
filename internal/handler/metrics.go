package handler

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// handleMetrics is a prometheus metrics handler.
func (s *handler) handleMetrics() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ReplaceAll(
			r.Header.Get("Authorization"),
			"Bearer ", "")

		if token != s.authToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := s.d.CollectMetrics(); err != nil {
			log.WithError(err).
				Error("failed to collect metrics from deCONZ")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		promhttp.Handler().ServeHTTP(w, r)
	}
}
