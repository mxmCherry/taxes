package taxes

import "time"

// CurrencyRates defines currency rates lookup.
type CurrencyRates interface {
	Rate(from, to string, date time.Time) (float64, error)
}
