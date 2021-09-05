package tax

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	Time     time.Time  `yaml:"time,omitempty"`
	Amount   *big.Float `yaml:"amount,omitempty"`
	Currency string     `yaml:"currency,omitempty"`

	CurrencyRate *big.Float `yaml:"currency_rate,omitempty"`
	LocalAmount  *big.Float `yaml:"local_amount,omitempty"`
	TaxRate      *big.Float `yaml:"tax_rate,omitempty"`
	TaxAmount    *big.Float `yaml:"tax_amount,omitempty"`
}

type Quarter struct {
	Quarter int `yaml:"quarter,omitempty"`

	// for this quarter only:
	TotalIncome *big.Float `yaml:"total_income,omitempty"`
	TotalTax    *big.Float `yaml:"total_tax,omitempty"`

	// since beginning of fiscal year:
	CumulativeIncome *big.Float `yaml:"cumulative_income,omitempty"`
	CumulativeTax    *big.Float `yaml:"cumulative_tax,omitempty"`

	Transactions []*Transaction `yaml:"transactions,omitempty"`
}

type Year struct {
	Year        int        `yaml:"year,omitempty"`
	TotalIncome *big.Float `yaml:"total_income,omitempty"`
	TotalTax    *big.Float `yaml:"total_tax,omitempty"`
	Quarters    []*Quarter `yaml:"quarters,omitempty"`
}

// ----------------------------------------------------------------------------

type CurrencyRates interface {
	Rate(ctx context.Context, date time.Time, from, to string) (*big.Float, error)
}

type Calc struct {
	CurrencyRates CurrencyRates `yaml:"-"`
	LocalCurrency string        `yaml:"local_currency,omitempty"`
	TaxRate       *big.Float    `yaml:"tax_rate,omitempty"`
}

func (c *Calc) Year(ctx context.Context, y *Year) error {
	y.TotalIncome = new(big.Float)
	y.TotalTax = new(big.Float)

	for _, q := range y.Quarters {
		q.CumulativeIncome = y.TotalIncome
		q.CumulativeTax = y.TotalTax

		if err := c.Quarter(ctx, q); err != nil {
			return fmt.Errorf("quarter %d: %w", q.Quarter, err)
		}

		y.TotalIncome = new(big.Float).Add(y.TotalIncome, q.TotalIncome)
		y.TotalTax = new(big.Float).Add(y.TotalTax, q.TotalTax)
	}

	return nil
}

func (c *Calc) Quarter(ctx context.Context, q *Quarter) error {
	q.TotalIncome = new(big.Float)
	q.TotalTax = new(big.Float)

	if q.CumulativeIncome == nil {
		q.CumulativeIncome = new(big.Float)
	}
	if q.CumulativeTax == nil {
		q.CumulativeTax = new(big.Float)
	}

	for _, tx := range q.Transactions {
		if err := c.Transaction(ctx, tx); err != nil {
			return fmt.Errorf("transaction %s: %w", tx.Time, err)
		}
		q.TotalIncome = new(big.Float).Add(q.TotalIncome, tx.LocalAmount)
		q.TotalTax = new(big.Float).Add(q.TotalTax, tx.TaxAmount)

		q.CumulativeIncome = new(big.Float).Add(q.CumulativeIncome, tx.LocalAmount)
		q.CumulativeTax = new(big.Float).Add(q.CumulativeTax, tx.TaxAmount)
	}

	return nil
}

func (c *Calc) Transaction(ctx context.Context, tx *Transaction) error {
	if tx.LocalAmount != nil {
		return nil
	}

	if tx.Currency == c.LocalCurrency {
		tx.LocalAmount = tx.Amount
	} else {
		if tx.CurrencyRate == nil {
			rate, err := c.CurrencyRates.Rate(ctx, tx.Time, tx.Currency, c.LocalCurrency)
			if err != nil {
				return fmt.Errorf("rate %q -> %q (%s): %w", tx.Currency, c.LocalCurrency, tx.Time, err)
			}
			tx.CurrencyRate = rate
		}
		tx.LocalAmount = new(big.Float).Mul(tx.Amount, tx.CurrencyRate)
	}

	if tx.TaxRate == nil {
		tx.TaxRate = c.TaxRate
	}

	if tx.TaxAmount == nil {
		tx.TaxAmount = new(big.Float).Mul(tx.LocalAmount, tx.TaxRate)
	}

	return nil
}

// ----------------------------------------------------------------------------

type CalcRun struct {
	Calc `yaml:",inline"`
	Data []*Year `yaml:"data,omitempty"`
}

func (r *CalcRun) Run(ctx context.Context) error {
	for _, y := range r.Data {
		if err := r.Calc.Year(ctx, y); err != nil {
			return fmt.Errorf("year %d: %w", y.Year, err)
		}
	}
	return nil
}
