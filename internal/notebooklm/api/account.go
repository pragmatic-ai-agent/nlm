package api

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/nlm/internal/notebooklm/rpc"
)

// AccountStatus is the current ZwVcOc account/status response.
//
// The wire response is not the generated Account proto. It is a compact status
// array with limit and feature data; older code decoded it positionally as an
// email address and settings fields.
type AccountStatus struct {
	NotebookLimit int `json:"notebook_limit,omitempty"`
	SourceLimit   int `json:"source_limit,omitempty"`
	UploadLimit   int `json:"upload_limit,omitempty"`
	Tier          int `json:"tier,omitempty"`
}

// GetAccountStatus dispatches ZwVcOc and decodes the current status shape.
func (c *Client) GetAccountStatus() (*AccountStatus, error) {
	resp, err := c.rpc.Do(rpc.Call{
		ID:   rpc.RPCGetOrCreateAccount,
		Args: []interface{}{},
	})
	if err != nil {
		return nil, fmt.Errorf("get account status: %w", err)
	}
	status, err := parseAccountStatus(resp)
	if err != nil {
		return nil, fmt.Errorf("get account status: decode response: %w", err)
	}
	return status, nil
}

func parseAccountStatus(b []byte) (*AccountStatus, error) {
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	items, ok := raw.([]any)
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("missing account status")
	}
	if len(items) == 1 {
		if nested, ok := items[0].([]any); ok {
			items = nested
		}
	}
	if len(items) < 2 {
		return nil, fmt.Errorf("missing account limits")
	}
	limits, ok := items[1].([]any)
	if !ok || len(limits) < 5 {
		return nil, fmt.Errorf("missing account limits")
	}
	values := make([]int, 5)
	for i := range values {
		n, ok := accountNumber(limits[i])
		if !ok {
			return nil, fmt.Errorf("bad account limit %d", i)
		}
		values[i] = int(n)
	}
	status := &AccountStatus{
		NotebookLimit: values[1],
		SourceLimit:   values[2],
		UploadLimit:   values[3],
		Tier:          values[4],
	}
	return status, nil
}

func accountNumber(v any) (float64, bool) {
	n, ok := v.(float64)
	return n, ok
}
