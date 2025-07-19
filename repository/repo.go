package repository

import "time"

type Metric struct {
	TraceID   string `db:"trace_id"`
	Timestamp string `json:"timestamp"`
	Stage     string `json:"stage"`
	Service   string `json:"service"`
}
type Raw struct {
	Timestamp time.Time
	TraceID   string
	Service   string
	Stage     uint8
}
type Stats struct {
	Service string
	Stage   int
	Count   int
}
