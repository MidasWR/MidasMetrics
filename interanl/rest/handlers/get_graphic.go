package handlers

import (
	"MidasMetrics/interanl/application/db_request"
	"MidasMetrics/interanl/rest/handlers/grafics"
	"errors"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog"
	"net/http"
)

func GetGraphicHandler(log zerolog.Logger, conn clickhouse.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		service := r.URL.Query().Get("service")
		if start == "" || end == "" || service == "" {
			log.Warn().Msgf("query parameters failed: start=%s end=%s service=%s", start, end, service)
			w.Write([]byte(errors.New("query parameters failed").Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics, err := db_request.GetMetrics(conn, service, start, end)
		if err != nil {
			log.Error().Msgf("get from db failed: %v", err)
			w.Write([]byte(errors.New("error select from db").Error()))
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}
		next := grafics.GraphicHandler(metrics, service, log)
		next.ServeHTTP(w, r)
	}
}
