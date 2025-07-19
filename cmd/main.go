package main

import (
	"MidasMetrics/config"
	"MidasMetrics/interanl/application/db_request/clickhouse"
	"MidasMetrics/interanl/rest"
	"MidasMetrics/repository"
	"context"
	"github.com/MidasWR/base-sdk-framework/midas"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	log := midas.InitLogger(midas.LoggerConfig{
		LogLevel: "PROD",
		Out:      os.Stdout,
	}, "Metrics")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	midas.RetryFunc[config.Config](
		midas.RetryConfig[config.Config]{
			Config:   *cfg,
			Logger:   log,
			Attempts: 5,
		},
		Start,
	)
}
func Start(cfg config.Config, logger zerolog.Logger) error {
	stch := repository.StorageCH{}
	stcs := repository.StorageCS{}
	stch.Init(cfg, logger)
	stcs.Init(cfg, logger)
	srv := rest.InitServer(cfg)
	go clickhouse.InsertMetricFunc(stch.Conn, logger, context.Background())
	srv.Run(logger, stch.Conn, stcs.Conn)
	return nil
}
