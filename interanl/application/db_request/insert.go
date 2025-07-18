package db_request

import (
	"MidasMetrics/repository"
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

type BatchInsert struct {
	Ch chan repository.Metric
}

var BI BatchInsert = BatchInsert{
	Ch: make(chan repository.Metric, 100000),
}

func InsertToChannel(metric repository.Metric) {
	BI.Ch <- metric
}
func InsertMetricFunc(conn clickhouse.Conn, log zerolog.Logger, ctx context.Context) {
	const query = `INSERT INTO metrics(trace_id,timestamp,service,stage)`
	batch, err := conn.PrepareBatch(ctx, query)
	if err != nil {
		log.Fatal().Err(err).Msg("initial prepare batch failed")
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case data := <-BI.Ch:
			parsedTime, err := time.Parse("02.01.2006 15:04", data.Timestamp)
			if err != nil {
				log.Error().Err(err).Msg("cant parsing timestamp")
				continue
			}

			stageVal, err := strconv.ParseUint(data.Stage, 10, 8)
			if err != nil {
				log.Error().Err(err).Msg("cant parsing stage")
				continue
			}
			fmt.Println(data.TraceID, data.Timestamp, data.Service, data.Stage)
			if err := batch.Append(data.TraceID, parsedTime, data.Service, uint8(stageVal)); err != nil {
				log.Error().Err(err).Msg("error appending to batch")
			}
		case <-ticker.C:
			if batch.Rows() > 0 {
				if err := batch.Send(); err != nil {
					log.Error().Err(err).Msg("error sending batch")
				}
				batch, err = conn.PrepareBatch(ctx, query)
				if err != nil {
					log.Error().Err(err).Msg("re-prepare batch failed")
					return
				}
			}

		case <-ctx.Done():
			if batch.Rows() > 0 {
				_ = batch.Send()
			}
			return
		}
	}
}
