package handlers

import (
	"MidasMetrics/interanl/application/db_request"
	"MidasMetrics/repository"
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

func PostMetric(log zerolog.Logger) http.HandlerFunc {
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
		db_request.InsertToChannel(Metric)
		w.WriteHeader(http.StatusCreated)
	}
}
