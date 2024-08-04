package bankgovua

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var errOnlyUAH = errors.New("only UAH target currency is supported")

type Client struct {
	HTTP interface {
		Do(*http.Request) (*http.Response, error)
	}
}

func (c *Client) Rate(ctx context.Context, at time.Time, from, to string) (float64, error) {
	if to != "UAH" {
		return 0, errOnlyUAH
	}

	req, err := http.NewRequestWithContext(ctx, "GET", BuildURL(at, from), nil)
	if err != nil {
		return 0, fmt.Errorf("prepare request: %w", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response body: %w", err)
	}

	return ParseResponse(body)
}
