package fakemetrics

import (
	"net/http"

	"time"

	"github.com/VictoriaMetrics/metrics"

	"github.com/wabel-v0/fakemetrics"
)

func main() {

	//  Configure metrics generation

	cfg := fakemetrics.Config{

		MetricPrefix: "app_",

		NumCounters: 5,

		NumGauges: 3,

		NumHistograms: 2,

		UpdateInterval: 2 * time.Second,

		Labels: map[string]string{

			"environment": "production",
		},
		// Set UpdateMetrics to true to periodically update generated metrics.
		// Set to false if you want to only register metrics  and not update them.
		UpdateMetrics: true,
	}

	// Create and start generator

	gen := fakemetrics.New(cfg)

	gen.Start()

	defer gen.Stop()

	// Set up metrics endpoint

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {

		metrics.WritePrometheus(w, true)

	})

	// Start server

	http.ListenAndServe(":8080", nil)

}
