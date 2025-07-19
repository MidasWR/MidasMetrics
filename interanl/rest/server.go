package rest

import (
	"MidasMetrics/config"
	"MidasMetrics/interanl/rest/handlers"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/MidasWR/base-sdk-framework/midas"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type Server struct {
	r   *mux.Router
	Cfg config.Config
}

func InitServer(cfg config.Config) *Server {
	return &Server{
		r:   mux.NewRouter(),
		Cfg: cfg,
	}
}
func (s *Server) Run(log zerolog.Logger, connCh clickhouse.Conn, connCs *gocql.Session) {
	midCfg := midas.MidConfig{
		Log:           log,
		TokenInHeader: false,
		HeaderIn:      true,
	}
	s.r.HandleFunc("/api/v1/metric", midas.Middleware(handlers.PostMetric(log, connCs), midCfg)).Methods("POST", "OPTIONS")
	s.r.HandleFunc("/api/v1/graphic", midas.Middleware(handlers.GetGraphicHandler(log, connCh), midCfg)).Methods("GET", "OPTIONS")
	s.r.HandleFunc("/api/v1/diagram", midas.Middleware(handlers.GetDiagramHandler(log, connCs), midCfg)).Methods("GET", "OPTIONS")
	log.Info().Msgf("rest server listening on port %s", s.Cfg.Port)
	log.Fatal().Err(http.ListenAndServe(s.Cfg.Host+":"+s.Cfg.Port, s.r)).Msgf("rest server die on port %v", s.Cfg.Port)
}
