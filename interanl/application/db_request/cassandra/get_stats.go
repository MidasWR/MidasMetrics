package cassandra

import (
	"MidasMetrics/repository"
	"github.com/gocql/gocql"
)

func GetStats(conn *gocql.Session, service string) ([]repository.Stats, error) {
	var stats []repository.Stats

	iter := conn.Query(`
		SELECT count, stage 
		FROM stage_counters 
		WHERE service = ?`,
		service,
	).Iter()

	var count int
	var stage int

	for iter.Scan(&count, &stage) {
		stats = append(stats, repository.Stats{
			Service: service,
			Stage:   stage,
			Count:   count,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return stats, nil
}
