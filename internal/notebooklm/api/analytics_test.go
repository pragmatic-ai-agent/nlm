package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProjectAnalyticsFixture(t *testing.T) {
	path := filepath.Join("testdata", "AUrzMb_analytics_response.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseProjectAnalytics(data)
	if err != nil {
		t.Fatalf("parseProjectAnalytics() error = %v", err)
	}
	if len(got.Series) != 2 {
		t.Fatalf("series count = %d, want 2", len(got.Series))
	}

	tests := []struct {
		i        int
		metricID int
		points   int
		last     int
	}{
		{0, 1, 30, 1},
		{1, 2, 30, 5},
	}
	for _, tt := range tests {
		series := got.Series[tt.i]
		if series.MetricID != tt.metricID {
			t.Errorf("series[%d].MetricID = %d, want %d", tt.i, series.MetricID, tt.metricID)
		}
		if len(series.Points) != tt.points {
			t.Errorf("series[%d] points = %d, want %d", tt.i, len(series.Points), tt.points)
			continue
		}
		last := series.Points[len(series.Points)-1]
		if last.Value != tt.last {
			t.Errorf("series[%d] last value = %d, want %d", tt.i, last.Value, tt.last)
		}
		if last.Time.Format("2006-01-02") != "2026-04-15" {
			t.Errorf("series[%d] last date = %s, want 2026-04-15", tt.i, last.Time.Format("2006-01-02"))
		}
	}
}
