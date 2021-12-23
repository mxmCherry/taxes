package tax_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = Describe("Calc", func() {
	It("calculates taxes", func() {
		subject, err := NewCalc(
			&Business{
				BaseCurrency:      "UAH",
				TaxRate:           0.05,
				RoundingPrecision: 2,
			},
			TransactionSlice{
				{Time: time.Date(2020, time.January, 1, 9, 45, 0, 0, time.UTC), Amount: 100.00, Currency: "GBP"},
				{Time: time.Date(2020, time.February, 2, 10, 15, 0, 0, time.UTC), Amount: 200.00, Currency: "USD"},
				{Time: time.Date(2020, time.September, 3, 14, 53, 0, 0, time.UTC), Amount: 300.00, Currency: "GBP"},
				{Time: time.Date(2021, time.October, 4, 13, 49, 0, 0, time.UTC), Amount: 400.00, Currency: "GBP"},
			},
			currencyRates{
				"20200101": {
					"GBP": {
						"UAH": 31.0206,
					},
				},
				"20200202": {
					"USD": {
						"UAH": 24.9196,
					},
				},
				"20200903": {
					"GBP": {
						"UAH": 36.9059,
					},
				},
				"20211004": {
					"GBP": {
						"UAH": 36.0313,
					},
				},
			},
		)
		Expect(err).NotTo(HaveOccurred())

		var res []*Quarter
		Expect(
			subject.Each(context.Background(), func(q *Quarter) error {
				res = append(res, q)
				return nil
			}),
		).To(Succeed())
		Expect(res).To(Equal([]*Quarter{
			// income: 100.00 * 31.0206 (3102.06) + 200.00 * 24.9196 (4983.92) = 8085.98
			// tax: 8085.98 * 0.05 = 404.299 ~ 404.30
			{Year: 2020, Quarter: 1, Income: 8085.98, Tax: 404.3},
			// income: 300.00 * 36.9059 (11071.77) + 8085.98 (Q1) = 19157.75
			// tax: 19157.75 * 0.05 = 957.8875 ~ 957.89
			{Year: 2020, Quarter: 3, Income: 19157.75, Tax: 957.89},
			// income: 400.00 * 36.0313 = 14412.52
			// tax: 14412.52 * 0.05 = 720.626 ~ 720.63
			{Year: 2021, Quarter: 4, Income: 14412.52, Tax: 720.63},
		}))
	})

	It("returns currency rates errors", func() {
		subject, err := NewCalc(
			&Business{
				BaseCurrency:      "UAH",
				TaxRate:           0.05,
				RoundingPrecision: 2,
			},
			TransactionSlice{
				{Time: time.Date(2020, time.January, 1, 9, 45, 0, 0, time.UTC), Amount: 100.00, Currency: "GBP"},
			},
			currencyRates{},
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(
			subject.Each(context.Background(), func(*Quarter) error { return nil }),
		).To(MatchError(HavePrefix(`get currency rate 2020-01-01 09:45:00 +0000 UTC "GBP" -> "UAH": `)))
	})

	It("rejects unsorted TransactionSlice", func() {
		subject, err := NewCalc(
			&Business{
				BaseCurrency:      "UAH",
				TaxRate:           0.05,
				RoundingPrecision: 2,
			},
			TransactionSlice{
				{Time: time.Date(2020, time.January, 1, 9, 45, 0, 1, time.UTC), Amount: 100.00, Currency: "GBP"}, // 1ns ahead
				{Time: time.Date(2020, time.January, 1, 9, 45, 0, 0, time.UTC), Amount: 100.00, Currency: "GBP"},
			},
			currencyRates{
				"20200101": {
					"GBP": {
						"UAH": 31.0206,
					},
				},
			},
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(
			subject.Each(context.Background(), func(*Quarter) error { return nil }),
		).To(MatchError(`transactions are not ordered by time, 2020-01-01 09:45:00 +0000 UTC should be before 2020-01-01 09:45:00.000000001 +0000 UTC`))
	})
})

type currencyRates map[string]map[string]map[string]float64

func (rr currencyRates) Rate(_ context.Context, at time.Time, from, to string) (float64, error) {
	date := at.Format("20060102")
	r, ok := rr[date][from][to]
	if !ok {
		return 0, fmt.Errorf("no rate for %s %q -> %q", date, from, to)
	}
	return r, nil
}
