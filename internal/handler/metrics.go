package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/davidborzek/deconz-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// authenticate authenticates a request when a token is configured.
func (s *handler) authenticate(r *http.Request) error {
	if len(s.expectedToken) == 0 {
		return nil
	}

	token := strings.ReplaceAll(
		r.Header.Get("Authorization"),
		"Bearer ", "")

	if token != s.expectedToken {
		return errors.New("authentication failed")
	}

	return nil
}

// handleMetrics is a prometheus metrics handler.
func (s *handler) handleMetrics() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.authenticate(r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := s.d.CollectMetrics(); err != nil {
			log.WithError(err).
				Error("failed to collect metrics from deCONZ")

			metrics.ErrorCounter.Inc()
		}

		promhttp.Handler().ServeHTTP(w, r)
	}
}
