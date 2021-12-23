package bankgovua

import (
	"encoding/json"
	"fmt"
	"time"
)

func BuildURL(at time.Time, from string) string {
	return fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json",
		from,
		at.Format("20060102"),
	)
}

func ParseResponse(p []byte) (float64, error) {
	var data []struct {
		Rate float64 `json:"rate"`
	}
	if err := json.Unmarshal(p, &data); err != nil {
		return 0, fmt.Errorf("decode JSON: %w (%s)", err, string(p))
	}
	if len(data) != 1 {
		return 0, fmt.Errorf("unexpected response size: %s", string(p))
	}
	if data[0].Rate == 0 {
		return 0, fmt.Errorf("no rate provided: %s", string(p))
	}
	return data[0].Rate, nil
}
