package taxes

import "time"

// Payment holds bank payment data.
type Payment struct {
	// Date holds payment date.
	Date time.Time
	// Currency holds payment currency, ISO 4217 (e.g., UAH).
	Currency string
	// Amount holds payment amount (in payment's currency).
	Amount float64
}
