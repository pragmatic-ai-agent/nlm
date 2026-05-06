package api

import (
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/tmc/nlm/gen/notebooklm/v1alpha1"
	intmethod "github.com/tmc/nlm/internal/method"
	"github.com/tmc/nlm/internal/notebooklm/rpc"
)

// ProjectAnalytics is the AUrzMb response: a set of metric time series.
type ProjectAnalytics struct {
	Series []AnalyticsSeries `json:"series"`
}

// AnalyticsSeries contains all buckets for one metric id.
type AnalyticsSeries struct {
	MetricID int              `json:"metric_id"`
	Points   []AnalyticsPoint `json:"points"`
}

// AnalyticsPoint is one dated metric bucket.
type AnalyticsPoint struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

// GetProjectAnalytics returns the AUrzMb time-series analytics for projectID.
func (c *Client) GetProjectAnalytics(projectID string) (*ProjectAnalytics, error) {
	req := &pb.GetProjectAnalyticsRequest{ProjectId: projectID}
	resp, err := c.rpc.Do(rpc.Call{
		ID:         rpc.RPCGetProjectAnalytics,
		NotebookID: projectID,
		Args:       intmethod.EncodeGetProjectAnalyticsArgsV2(req),
	})
	if err != nil {
		return nil, fmt.Errorf("get project analytics: %w", err)
	}
	analytics, err := parseProjectAnalytics(resp)
	if err != nil {
		return nil, fmt.Errorf("get project analytics: decode response: %w", err)
	}
	return analytics, nil
}

func parseProjectAnalytics(b []byte) (*ProjectAnalytics, error) {
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	items, ok := analyticsSeriesItems(raw)
	if !ok {
		return nil, fmt.Errorf("missing metric series")
	}
	var out ProjectAnalytics
	for _, item := range items {
		series, err := parseAnalyticsSeries(item)
		if err != nil {
			return nil, err
		}
		out.Series = append(out.Series, series)
	}
	return &out, nil
}

func analyticsSeriesItems(v any) ([]any, bool) {
	a, ok := v.([]any)
	if !ok {
		return nil, false
	}
	if len(a) == 0 {
		return nil, true
	}
	if isAnalyticsSeries(a[0]) {
		return a, true
	}
	if len(a) == 1 {
		return analyticsSeriesItems(a[0])
	}
	return nil, false
}

func isAnalyticsSeries(v any) bool {
	a, ok := v.([]any)
	if !ok || len(a) < 2 {
		return false
	}
	_, ok = number(a[0])
	return ok
}

func parseAnalyticsSeries(v any) (AnalyticsSeries, error) {
	a, ok := v.([]any)
	if !ok || len(a) < 2 {
		return AnalyticsSeries{}, fmt.Errorf("bad metric series")
	}
	id, ok := number(a[0])
	if !ok {
		return AnalyticsSeries{}, fmt.Errorf("bad metric id")
	}
	buckets, ok := analyticsBuckets(a[1])
	if !ok {
		return AnalyticsSeries{}, fmt.Errorf("metric %d: missing buckets", int(id))
	}
	series := AnalyticsSeries{MetricID: int(id)}
	for _, bucket := range buckets {
		point, err := parseAnalyticsPoint(bucket)
		if err != nil {
			return AnalyticsSeries{}, fmt.Errorf("metric %d: %w", series.MetricID, err)
		}
		series.Points = append(series.Points, point)
	}
	return series, nil
}

func analyticsBuckets(v any) ([]any, bool) {
	a, ok := v.([]any)
	if !ok {
		return nil, false
	}
	if len(a) == 0 {
		return nil, true
	}
	if isAnalyticsBucket(a[0]) {
		return a, true
	}
	if len(a) == 1 {
		return analyticsBuckets(a[0])
	}
	for _, item := range a {
		if buckets, ok := analyticsBuckets(item); ok {
			return buckets, true
		}
	}
	return nil, false
}

func isAnalyticsBucket(v any) bool {
	a, ok := v.([]any)
	if !ok || len(a) < 3 {
		return false
	}
	ts, ok := a[0].([]any)
	if !ok || len(ts) == 0 {
		return false
	}
	if _, ok := number(ts[0]); !ok {
		return false
	}
	_, ok = number(a[2])
	return ok
}

func parseAnalyticsPoint(v any) (AnalyticsPoint, error) {
	a, ok := v.([]any)
	if !ok || len(a) < 3 {
		return AnalyticsPoint{}, fmt.Errorf("bad bucket")
	}
	ts, ok := a[0].([]any)
	if !ok || len(ts) == 0 {
		return AnalyticsPoint{}, fmt.Errorf("bad timestamp")
	}
	sec, ok := number(ts[0])
	if !ok {
		return AnalyticsPoint{}, fmt.Errorf("bad timestamp")
	}
	value, ok := number(a[2])
	if !ok {
		return AnalyticsPoint{}, fmt.Errorf("bad value")
	}
	return AnalyticsPoint{
		Time:  time.Unix(int64(sec), 0).UTC(),
		Value: int(value),
	}, nil
}

func number(v any) (float64, bool) {
	n, ok := v.(float64)
	return n, ok
}
