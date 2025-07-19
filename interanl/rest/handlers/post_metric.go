package handlers

import (
	"MidasMetrics/interanl/application/db_request/cassandra"
	"MidasMetrics/interanl/application/db_request/clickhouse"
	"MidasMetrics/repository"
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"net/http"
)

func PostMetric(log zerolog.Logger, conn *gocql.Session) http.HandlerFunc {
	var Metric repository.Metric
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&Metric)
		if err != nil {
			log.Error().Err(err).Msg("error decoding metric")
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := Validate(Metric); err != nil {
			log.Error().Err(err).Msg("error validating metric")
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clickhouse.InsertToChannel(Metric)
		if err = cassandra.UpdateCounter(conn, Metric); err != nil {
			log.Error().Err(err).Msg("error updating metric")
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
