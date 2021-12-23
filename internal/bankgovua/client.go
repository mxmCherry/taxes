package bankgovua

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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

	resp, err := c.get(ctx, BuildURL(at, from))
	if err != nil {
		return 0, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response body: %w", err)
	}

	return ParseResponse(body)
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	return c.HTTP.Do(req)
}
