package handlers

import (
	"MidasMetrics/repository"
	"errors"
)

func Validate(metric repository.Metric) error {
	if metric.Stage == "" {
		return errors.New("stage must be set")
	}
	if metric.Service == "" {
		return errors.New("service must be set")
	}
	return nil
}
