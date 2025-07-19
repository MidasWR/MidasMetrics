package cassandra

import (
	"MidasMetrics/repository"
	"github.com/gocql/gocql"
	"strconv"
)

func UpdateCounter(conn *gocql.Session, metric repository.Metric) error {
	stage, err := strconv.Atoi(metric.Stage)
	if err != nil {
		return err
	}
	err = conn.Query(`UPDATE stage_counters SET count = count + 1 WHERE service= ? AND stage= ?;
;`, metric.Service, stage).Exec()
	if err != nil {
		return err
	}
	return nil
}
