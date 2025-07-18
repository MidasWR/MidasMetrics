package grafics

import (
	"MidasMetrics/repository"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog"

	"net/http"
)

func GraphicHandler(metric []repository.Raw, service string, log zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				AssetsHost: "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/",
			}),
			charts.WithTitleOpts(opts.Title{
				Title: "Tracing: " + service,
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: opts.Bool(true),
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Type:      "category",
				Name:      "Time",
				AxisLabel: &opts.AxisLabel{Show: opts.Bool(false)},
				SplitLine: &opts.SplitLine{Show: opts.Bool(false)},
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "Stage",
				Type: "value",
			}),
		)

		var xAxis []string
		var yAxis []opts.LineData
		for _, m := range metric {
			xAxis = append(xAxis, m.Timestamp.Format("2006-01-02 15:04:05"))
			stageInt := int(m.Stage)
			yAxis = append(yAxis, opts.LineData{Value: stageInt})
		}

		line.SetXAxis(xAxis).AddSeries("Stage", yAxis).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))
		w.Header().Set("Content-Type", "text/html")
		err := line.Render(w)
		if err != nil {
			log.Error().Err(err).Msg("error rendering graphics line")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
