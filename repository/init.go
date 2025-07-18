package repository

import (
	"MidasMetrics/config"
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog"
	"net"
	"time"
)

type Storage struct {
	Conn clickhouse.Conn
}

func (s *Storage) InitStorage(cfg config.Config, log zerolog.Logger) {
	var addr []string
	var err error
	addr = append(addr, cfg.DBHost+":"+cfg.DBPort)
	s.Conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: addr,
		Auth: clickhouse.Auth{
			Database: cfg.DBHandler,
			Username: "default",
			Password: cfg.DBPassword,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msgf("error connecting to database on port %s", cfg.DBPort)
	}
	defer func() {
		if err := s.Conn.Close(); err != nil {
			log.Error().Err(err).Msg("error closing ClickHouse connection")
		}
	}()
	log.Info().Msgf("connected to database on port %s", cfg.DBPort)
	q := `
		CREATE TABLE IF NOT EXISTS metrics (
		trace_id String,
		timestamp DateTime,
		service String,
		stage UInt8,
		) ENGINE = ReplacingMergeTree(timestamp)
		ORDER BY (trace_id)`
	if err := s.Conn.Exec(context.Background(), q); err != nil {
		log.Error().Err(err).Msg("error creating table")
		return
	}
	defer s.Conn.Close()
	log.Info().Msg("table 'metrics' ensured")
}
