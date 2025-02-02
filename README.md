
# FakeMetrics

FakeMetrics is a Go package designed for generating fake metrics that can be integrated with Prometheus or Victoria Metrics. This package is ideal for testing and demo purposes, providing a simple way to simulate metric data without needing a live system.

## Features

- **Fake Counters**: Automatically increments counters with random values.
- **Fake Gauges**: Generates gauges that return random float values.
- **Fake Histograms**: Updates histograms with random values.
- **Customizable**: Configure the number of each metric type, update intervals, prefixes, and static labels.
- **Update Metrics Toggle**: Choose whether to periodically update generated metrics via the `UpdateMetrics` configuration flag.

## Installation

To use FakeMetrics in your Go project, install it using `go get`:

bash

```bash
go get github.com/wabel-v0/fakemetrics
```

Then import it in your code:

go

```go
import (
    "github.com/wabel-v0/fakemetrics"
)
```

## Usage

Below is a basic example of how to use FakeMetrics:

go

```go
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
// Set UpdateMetrics to true to update generated metrics.
// Set to false if you want to only register metrics and not update them.
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
```

## Configuration

The `Config` struct allows you to customize the behavior of the metrics generator:

go

```go
type Config struct {
    MetricPrefix   string            // Prefix added to each metric name (default: "fake_")
    NumCounters    int               // Number of counters to create (default: 10)
    NumGauges      int               // Number of gauges to create (default: 10)
    NumHistograms  int               // Number of histograms to create (default: 10)
    UpdateInterval time.Duration     // Interval for updating metric values (default: 2s)
    Labels         map[string]string // A set of labels to be appended to each metric (default: {"environment": "lazy"})
	UpdateMetrics  bool              // If true, it update generated metrics (default: false)

}
```

## How It Works

1. **Metric Creation**: Upon starting, the generator creates the specified number of counters, gauges, and histograms. Each metric name is prefixed with `MetricPrefix`, and additional labels are appended in Prometheus style (e.g., `key="value"`).
    
2. **Metric Updates**: If enabled by the `UpdateMetrics` flag, a goroutine runs on a ticker based on the `UpdateInterval`. Each tick performs the following:
- Counters are incremented by a random value.
- Gauges provide a new random float value (via a dynamic function).
- Histograms are updated with new random float values.

    

## Dependencies

This project uses the following libraries:

- [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics): For metric registration and manipulation.

Feel free to contribute or open issues if you have any questions. Enjoy testing!

