package fakemetrics

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/stretchr/testify/assert"
)

func uniquePrefix(base string) string {
	return fmt.Sprintf("%s_%d_", base, time.Now().UnixNano())
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name:   "default configuration",
			config: Config{},
			expected: Config{
				MetricPrefix:   "app_",
				NumCounters:    10,
				NumGauges:      10,
				NumHistograms:  10,
				UpdateInterval: 2 * time.Second,
				Labels: map[string]string{
					"environment": "lazy",
				},
			},
		},
		{
			name: "custom configuration",
			config: Config{
				MetricPrefix:   "custom_",
				NumCounters:    5,
				NumGauges:      5,
				NumHistograms:  5,
				UpdateInterval: 1 * time.Second,
				Labels: map[string]string{
					"env": "test",
				},
			},
			expected: Config{
				MetricPrefix:   "custom_",
				NumCounters:    5,
				NumGauges:      5,
				NumHistograms:  5,
				UpdateInterval: 1 * time.Second,
				Labels: map[string]string{
					"env": "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := New(tt.config)
			assert.Equal(t, tt.expected, generator.config)
			assert.NotNil(t, generator.stopChan)
		})
	}
}

func TestBuildName(t *testing.T) {
	tests := []struct {
		name       string
		config     Config
		metricName string
		expected   string
	}{
		{
			name: "with labels",
			config: Config{
				MetricPrefix: "test_",
				Labels: map[string]string{
					"env":  "prod",
					"zone": "us-east",
				},
			},
			metricName: "metric_1",
			expected:   `test_metric_1{env="prod",zone="us-east"}`,
		},
		{
			name: "without labels",
			config: Config{
				MetricPrefix: "test_",
				Labels:       map[string]string{},
			},
			metricName: "metric_1",
			expected:   "test_metric_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := New(tt.config)
			result := generator.buildName(tt.metricName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStartStop(t *testing.T) {
	// Generate a unique prefix for this test to avoid duplicate registrations.
	prefix := uniquePrefix("testStartStop")
	config := Config{
		MetricPrefix:   prefix,
		NumCounters:    2,
		NumGauges:      2,
		NumHistograms:  2,
		UpdateInterval: 100 * time.Millisecond,
		Labels: map[string]string{
			"env": "test",
		},
	}

	generator := New(config)
	generator.Start()

	for i := 0; i < 5; i++ {
		generator.updateMetrics()
		time.Sleep(50 * time.Millisecond)
	}

	// Extra wait for any asynchronous processing.
	time.Sleep(200 * time.Millisecond)

	var buf bytes.Buffer
	metrics.WritePrometheus(&buf, false)
	metricsString := buf.String()

	// Build expected names for counters
	expectedCounter0 := fmt.Sprintf(`%scounter_0{env="test"}`, prefix)
	expectedCounter1 := fmt.Sprintf(`%scounter_1{env="test"}`, prefix)
	// For histograms, check the _count series which is always exported.
	expectedHistogram0 := fmt.Sprintf(`%shistogram_0_count{env="test"}`, prefix)
	expectedHistogram1 := fmt.Sprintf(`%shistogram_1_count{env="test"}`, prefix)

	assert.True(t, strings.Contains(metricsString, expectedCounter0),
		"Counter 0 %q not found in output", expectedCounter0)
	assert.True(t, strings.Contains(metricsString, expectedCounter1),
		"Counter 1 %q not found in output", expectedCounter1)
	assert.True(t, strings.Contains(metricsString, expectedHistogram0),
		"Histogram 0 count %q not found in output", expectedHistogram0)
	assert.True(t, strings.Contains(metricsString, expectedHistogram1),
		"Histogram 1 count %q not found in output", expectedHistogram1)

	generator.Stop()
}

func TestCreateMetrics(t *testing.T) {
	// Generate a unique prefix for this test.
	prefix := uniquePrefix("testCreateMetrics")
	config := Config{
		MetricPrefix:  prefix,
		NumCounters:   2,
		NumGauges:     2,
		NumHistograms: 2,
		Labels: map[string]string{
			"env": "test",
		},
	}

	generator := New(config)
	// Create the metrics.
	generator.createMetrics()

	for i := 0; i < 5; i++ {
		generator.updateMetrics()
		time.Sleep(50 * time.Millisecond)
	}

	// Extra wait to allow updates to be flushed.
	time.Sleep(200 * time.Millisecond)

	var buf bytes.Buffer
	metrics.WritePrometheus(&buf, false)
	metricsString := buf.String()

	// Verify counters and gauges are created.
	for i := 0; i < config.NumCounters; i++ {
		metricName := fmt.Sprintf(`%scounter_%d{env="test"}`, prefix, i)
		assert.True(t, strings.Contains(metricsString, metricName),
			"Expected counter %d with name %q not found", i, metricName)
	}

	for i := 0; i < config.NumGauges; i++ {
		metricName := fmt.Sprintf(`%sgauge_%d{env="test"}`, prefix, i)
		assert.True(t, strings.Contains(metricsString, metricName),
			"Expected gauge %d with name %q not found", i, metricName)
	}

	// For histograms, check the _count series.
	for i := 0; i < config.NumHistograms; i++ {
		metricName :=
			fmt.Sprintf(`%shistogram_%d_count{env="test"}`, prefix, i)
		assert.True(t, strings.Contains(metricsString, metricName),
			"Expected histogram %d count with name %q not found", i, metricName)
	}
}
