package repository

import (
	"MidasMetrics/config"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

type StorageCS struct {
	Conn *gocql.Session
}

func (s *StorageCS) Init(cfg config.Config, log zerolog.Logger) {
	cluster := gocql.NewCluster(cfg.Host)
	cluster.Keyspace = "midas_metrics"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 60 * time.Second
	cluster.SerialConsistency = gocql.LocalSerial
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	port, _ := strconv.Atoi(cfg.CSPort)
	cluster.Port = port
	cluster.NumConns = 5
	var err error
	s.Conn, err = cluster.CreateSession()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to cassandra")
	}
	err = s.Conn.Query(`CREATE TABLE IF NOT EXISTS stage_counters (    service text,    stage int,    count counter,   PRIMARY KEY (service, stage));`).Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("error creating stage counters")
	}
	go func() {
		ticker := time.NewTicker(24 * time.Hour)

		for {
			select {
			case <-ticker.C:
				if err := s.Conn.Query(`TRUNCATE TABLE stage_counters`).Exec(); err != nil {
					log.Panic().Err(err).Msg("error truncating stage counters")
				}
			}
		}

	}()
	log.Info().Msg("stage counters table created")
}
