package handlers

import (
	"MidasMetrics/interanl/application/db_request/cassandra"
	"MidasMetrics/interanl/rest/handlers/view"
	"errors"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"net/http"
)

func GetDiagramHandler(log zerolog.Logger, conn *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")

		stats, err := cassandra.GetStats(conn, service)
		if err != nil {
			if errors.Is(err, gocql.ErrNotFound) {
				log.Warn().Msgf("error found stats for %s", service)
				w.Write([]byte("error found stats for " + service))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error().Err(err).Msgf("error getting stats for %s", service)
			w.Write([]byte("error getting stats for " + service))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		view.DiagramHandler(stats, service, log).ServeHTTP(w, r)

	}
}
