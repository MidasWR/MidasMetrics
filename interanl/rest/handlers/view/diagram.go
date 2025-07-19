package view

import (
	"MidasMetrics/repository"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog"
	"net/http"
)

func DiagramHandler(stats []repository.Stats, service string, log zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		bar := charts.NewBar()
		bar.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title:    "Stage Counters: " + service,
				Subtitle: "Total stage counts",
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
			charts.WithXAxisOpts(opts.XAxis{Name: "Stage"}),
			charts.WithYAxisOpts(opts.YAxis{Name: "Count"}),
		)

		var stages []string
		var counts []opts.BarData

		for _, s := range stats {
			stages = append(stages, fmt.Sprintf("Stage %d", s.Stage))
			counts = append(counts, opts.BarData{Value: s.Count})
		}

		bar.SetXAxis(stages).
			AddSeries("Count", counts).
			SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: opts.Bool(true)}))

		w.Header().Set("Content-Type", "text/html")
		if err := bar.Render(w); err != nil {
			log.Error().Err(err).Msg("failed to render bar chart")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
