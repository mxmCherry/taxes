package bankgovua

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// DefaultApiAddr holds default API hostname/domain.
const DefaultApiAddr = "https://bank.gov.ua"

var (
	errUnsupportedTo = errors.New("bankgovua: unsupported 'to' currency, only UAH available")
	errUnexpectedRes = errors.New("bankgovua: unexpected result")
)

// CurrencyRates implemets currency rates client for bank.gov.ua:
// https://bank.gov.ua/control/uk/publish/article?art_id=38441973
type CurrencyRates struct {
	// http holds HTTP client.
	http *http.Client
	// addr holds API hostname/domain.
	addr string
	// cache holds currency rates cache "{from}-{to}-{YYYYMMDD}" -> RATE.
	cache map[string]float64
}

// NewCurrencyRates constructs new currency rates client for bank.gov.ua.
func NewCurrencyRates() *CurrencyRates {
	return NewCurrencyRatesWithAddr(DefaultApiAddr)
}

// NewCurrencyRatesWithHost constructs new currency rates client for given addr (like proto://host:port, no path).
// Mostly for testing.
func NewCurrencyRatesWithAddr(addr string) *CurrencyRates {
	return &CurrencyRates{
		http:  http.DefaultClient,
		addr:  addr,
		cache: map[string]float64{},
	}
}

// Rate returns currency rate for given date.
// Only to=UAH is supported.
func (r *CurrencyRates) Rate(from, to string, date time.Time) (float64, error) {
	if from == to {
		return 1, nil
	}
	if to != "UAH" {
		return 0, errUnsupportedTo
	}

	dateStr := date.Format("20060102")
	cacheKey := fmt.Sprintf("%s-%s-%s", from, to, dateStr)
	if rate, ok := r.cache[cacheKey]; ok {
		return rate, nil
	}

	resp, err := r.http.Get(fmt.Sprintf(
		"%s/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json",
		r.addr,
		from,
		dateStr,
	))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var res []struct {
		Rate float64 `json:"rate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, err
	}
	if len(res) != 1 {
		return 0, errUnexpectedRes
	}

	rate := res[0].Rate
	r.cache[cacheKey] = rate
	return rate, nil
}
