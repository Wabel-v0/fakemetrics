package fakemetrics

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"golang.org/x/exp/rand"
)

// Config defines the fake metrics generation parameters
type Config struct {
	MetricPrefix   string
	NumCounters    int
	NumGauges      int
	NumHistograms  int
	UpdateInterval time.Duration
	Labels         map[string]string
	UpdateMetrics  bool
}

// Generator creates and manages fake metrics
type Generator struct {
	config    Config
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
}

// New creates a metrics generator with default values
func New(cfg Config) *Generator {
	// Set default values
	if cfg.MetricPrefix == "" {
		cfg.MetricPrefix = "fake_"
	}
	if cfg.NumCounters == 0 {
		cfg.NumCounters = 10
	}
	if cfg.NumGauges == 0 {
		cfg.NumGauges = 10
	}
	if cfg.NumHistograms == 0 {
		cfg.NumHistograms = 10
	}
	if cfg.UpdateInterval == 0 {
		cfg.UpdateInterval = 2 * time.Second
	}
	if cfg.Labels == nil {
		cfg.Labels = map[string]string{
			"environment": "lazy",
		}
	}

	return &Generator{
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

// Start begins updating metric values
func (g *Generator) Start() {
	g.createMetrics()

	if g.config.UpdateMetrics {
		g.waitGroup.Add(1)
		go func() {
			defer g.waitGroup.Done()
			ticker := time.NewTicker(g.config.UpdateInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					g.updateMetrics()
				case <-g.stopChan:
					return
				}
			}
		}()
	} else {
		fmt.Println("Metric updates are disabled by configuration.")
	}
}

// Stop halts metric updates
func (g *Generator) Stop() {
	close(g.stopChan)
	g.waitGroup.Wait()
}

func (g *Generator) createMetrics() {
	// Create counters
	for i := 0; i < g.config.NumCounters; i++ {
		name := g.buildName(fmt.Sprintf("counter_%d", i))
		metrics.NewCounter(name)
	}

	// Create gauges
	for i := 0; i < g.config.NumGauges; i++ {
		name := g.buildName(fmt.Sprintf("gauge_%d", i))
		metrics.NewGauge(name, func() float64 {
			return rand.Float64() * 100
		})
	}

	// Create histograms
	for i := 0; i < g.config.NumHistograms; i++ {
		name := g.buildName(fmt.Sprintf("histogram_%d", i))
		metrics.NewHistogram(name)
	}
}

func (g *Generator) updateMetrics() {
	// Update counter values using valid range
	for i := 0; i < g.config.NumCounters; i++ {
		name := g.buildName(fmt.Sprintf("counter_%d", i))
		intValue := rand.Intn(10) + 1
		metrics.GetOrCreateCounter(name).Add(intValue)

	}

	// Update histogram values
	for i := 0; i < g.config.NumHistograms; i++ {
		name := g.buildName(fmt.Sprintf("histogram_%d", i))
		metrics.GetOrCreateHistogram(name).Update(rand.Float64() * 100)
	}
}

func (g *Generator) buildName(metricName string) string {
	var labels []string
	for k, v := range g.config.Labels {
		labels = append(labels, fmt.Sprintf(`%s="%s"`, k, v))
	}

	if len(labels) == 0 {
		return g.config.MetricPrefix + metricName
	}

	return fmt.Sprintf("%s%s{%s}", g.config.MetricPrefix, metricName, strings.Join(labels, ","))
}
