package httpclient

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	client      *http.Client
	rateLimiter *rate.Limiter
	maxTries    int
}

func New(oneRequestPer time.Duration, maxTries int) *Client {
	return &Client{
		client:      new(http.Client),
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
		maxTries:    maxTries,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for i := 0; i < c.maxTries; i++ {
		c.rateLimiter.Wait(req.Context())

		resp, err = c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}
		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			continue
		}
	}
	return resp, nil
}
