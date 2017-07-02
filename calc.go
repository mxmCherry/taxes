package taxes

import (
	"sort"
)

// Calc is a tax calculator.
type Calc struct {
	rates CurrencyRates
}

// NewCalc constructs new quarterly tax calculator.
func NewCalc(rates CurrencyRates) *Calc {
	return &Calc{
		rates: rates,
	}
}

// Calc calculates company taxes.
func (c *Calc) Calc(company Company, payments []Payment) ([]Tax, error) {
	byQuarter := map[Quarter]Tax{}
	for _, p := range payments {
		bcp, err := CalcBaseCurrencyPayment(c.rates, company.BaseCurrency, p)
		if err != nil {
			return nil, err
		}

		q := QuarterOf(bcp.Payment.Date)
		t := byQuarter[q]
		t.Quarter = q
		t.Payments = append(t.Payments, *bcp)
		t.Income += bcp.Amount
		t.Tax = t.Income * company.TaxRate
		byQuarter[q] = t
	}

	taxes := make([]Tax, 0, len(byQuarter))
	for _, t := range byQuarter {
		sort.Sort(baseCurrencyPaymentsByDate(t.Payments))
		taxes = append(taxes, t)
	}

	sort.Sort(taxesByQuarter(taxes))

	return taxes, nil
}
