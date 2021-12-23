package bankgovua

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var errOnlyUAH = errors.New("only UAH target currency is supported")

type Client struct {
	HTTP        *http.Client
	RateLimiter *rate.Limiter // rate.Every(time.Second) - allow one RPS
	MaxRetries  int           // per one Rate call
}

func (c *Client) Rate(ctx context.Context, date time.Time, from, to string) (float64, error) {
	if to != "UAH" {
		return 0, errOnlyUAH
	}

	url := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json",
		from,
		date.Format("20060102"),
	)
	resp, err := c.get(ctx, url, c.MaxRetries+1)
	if err != nil {
		return 0, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response body: %w", err)
	}

	var data []struct {
		Rate float64 `json:"rate"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("decode JSON: %w (%s)", err, string(body))
	}
	if len(data) != 1 {
		return 0, fmt.Errorf("unexpected response size: %s", string(body))
	}
	if data[0].Rate == 0 {
		return 0, fmt.Errorf("no rate provided: %s", string(body))
	}
	return data[0].Rate, nil
}

func (c *Client) get(ctx context.Context, url string, maxTries int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	c.RateLimiter.Wait(ctx)

	var resp *http.Response
	for i := 0; i < maxTries; i++ {
		resp, err = c.HTTP.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			continue
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return nil, fmt.Errorf("unexpected response status %s: %s", resp.Status, string(body))
}
