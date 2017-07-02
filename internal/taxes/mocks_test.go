package taxes_test

import "time"

type mockCurrencyRates struct {
	request struct {
		from string
		to   string
		date time.Time
	}
	rate float64
	err  error
}

func (r *mockCurrencyRates) Rate(from, to string, date time.Time) (float64, error) {
	r.request.from = from
	r.request.to = to
	r.request.date = date

	return r.rate, r.err
}
