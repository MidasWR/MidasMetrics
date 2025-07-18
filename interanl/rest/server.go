package rest

import (
	"MidasMetrics/config"
	"MidasMetrics/interanl/rest/handlers"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/MidasWR/base-sdk-framework/midas"
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
func (s *Server) Run(log zerolog.Logger, conn clickhouse.Conn) {
	midCfg := midas.MidConfig{
		Log:           log,
		TokenInHeader: false,
		HeaderIn:      true,
	}
	s.r.HandleFunc("/api/v1/metric", midas.Middleware(handlers.PostMetric(log), midCfg)).Methods("POST", "OPTIONS")
	s.r.HandleFunc("/api/v1/grafic", midas.Middleware(handlers.GetGraphicHandler(log, conn), midCfg)).Methods("GET", "OPTIONS")
	log.Info().Msgf("rest server listening on port %s", s.Cfg.Port)
	log.Fatal().Err(http.ListenAndServe(s.Cfg.Host+":"+s.Cfg.Port, s.r)).Msgf("rest server die on port %v", s.Cfg.Port)
}
