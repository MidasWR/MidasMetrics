package db_request

import (
	"MidasMetrics/repository"
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"
)

func GetMetrics(conn clickhouse.Conn, service, start, end string) ([]repository.Raw, error) {
	defer conn.Close()
	q := `SELECT * FROM metrics FINAL WHERE timestamp >= ? AND timestamp <= ? AND service = ? ORDER BY timestamp`
	startParsed, err := ParseClickhouseTime(start)
	if err != nil {
		return nil, err
	}
	endParsed, err := ParseClickhouseTime(end)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), q, startParsed, endParsed, service)
	if err != nil {
		return nil, err
	}

	var metrics []repository.Raw
	for rows.Next() {
		var metric repository.Raw
		if err := rows.Scan(&metric.TraceID, &metric.Timestamp, &metric.Service, &metric.Stage); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics found")
	}
	return metrics, nil
}
func ParseClickhouseTime(input string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"02.01.2006 15:04",
		"02.01.2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cant handle time: %s", input)
}
