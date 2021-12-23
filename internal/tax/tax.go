package tax

import (
	"context"
	"fmt"
	"regexp"
	"time"
)

var currencyCodeRx = regexp.MustCompile(`^[A-Z]{3}$`)

type Business struct {
	BaseCurrency      string  `yaml:"base_currency"`
	TaxRate           float64 `yaml:"tax_rate"`
	RoundingPrecision int     `yaml:"rounding_precision,omitempty"`
}

func (b *Business) Validate() error {
	if b == nil {
		return fmt.Errorf("cannot be nil")
	}
	if !currencyCodeRx.MatchString(b.BaseCurrency) {
		return fmt.Errorf(`base_currency %q is invalid: should be a 3-letter ISO code, uppercased (like "UAH")`, b.BaseCurrency)
	}
	if b.TaxRate <= 0 || b.TaxRate >= 1 {
		return fmt.Errorf("tax_rate %f is invalid: should be in 0..1 range", b.TaxRate)
	}
	if b.RoundingPrecision <= 0 {
		return fmt.Errorf("rounding_precision %d is invalid: should be greater than 0", b.RoundingPrecision)
	}
	return nil
}

// ----------------------------------------------------------------------------

type Transactions interface {
	Each(func(*Transaction) error) error
}

type Transaction struct {
	Time     time.Time `yaml:"time"`
	Amount   float64   `yaml:"amount"`
	Currency string    `yaml:"currency"`
}

func (t *Transaction) Validate() error {
	if t == nil {
		return fmt.Errorf("cannot be nil")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("amount %f if invalid: should be greater than 0", t.Amount)
	}
	if !currencyCodeRx.MatchString(t.Currency) {
		return fmt.Errorf(`currency %q is invalid: should be a 3-letter ISO code, uppercased (like "UAH")`, t.Currency)
	}
	return nil
}

type TransactionSlice []*Transaction

func (tt TransactionSlice) Each(cb func(*Transaction) error) error {
	for _, t := range tt {
		if err := cb(t); err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------------------

type CurrencyRates interface {
	Rate(ctx context.Context, at time.Time, from, to string) (float64, error)
}

// ----------------------------------------------------------------------------

type Calc interface {
	Each(context.Context, func(*Quarter) error) error
}

type Quarter struct {
	Year    int
	Quarter int
	Income  float64
	Tax     float64
}

// ----------------------------------------------------------------------------

type calc struct {
	business     *Business
	transactions Transactions
	rates        CurrencyRates
}

func NewCalc(business *Business, transactions Transactions, rates CurrencyRates) (Calc, error) {
	if err := business.Validate(); err != nil {
		return nil, fmt.Errorf("business is invalid: %w", err)
	}

	return &calc{
		business:     business,
		transactions: transactions,
		rates:        rates,
	}, nil
}

func (c *calc) Each(ctx context.Context, cb func(*Quarter) error) error {
	var (
		prev *Transaction
		q    *Quarter
	)

	emit := func() error {
		q.Income = Round(q.Income, c.business.RoundingPrecision)
		q.Tax = Round(q.Income*c.business.TaxRate, c.business.RoundingPrecision)
		if err := cb(q); err != nil {
			return err
		}
		return nil
	}

	err := c.transactions.Each(func(tx *Transaction) error {
		// sanity checks
		if err := tx.Validate(); err != nil {
			return fmt.Errorf("transaction %v is invalid: %w", *tx, err)
		}
		if prev != nil && tx.Time.Before(prev.Time) {
			return fmt.Errorf("transactions are not ordered by time, %s should be before %s", tx.Time, prev.Time)
		}
		prev = tx

		txYear := tx.Time.Year()
		txQuarter := QuarterOf(tx.Time.Month())

		if q != nil {
			if txYear != q.Year {
				// end of year - emit accumulated quarter and init fresh one
				emit()
				q = &Quarter{
					Year:    txYear,
					Quarter: txQuarter,
					Income:  0,
					Tax:     0,
				}
			} else if txQuarter != q.Quarter {
				// same year, but next quarter - emit accumulated quarter and init fresh one with cumulative income/tax
				emit()
				q = &Quarter{
					Year:    txYear,
					Quarter: txQuarter,
					Income:  q.Income,
					Tax:     q.Tax,
				}
			}
		} else {
			// init quarter for the first time
			q = &Quarter{
				Year:    txYear,
				Quarter: txQuarter,
				Income:  0,
				Tax:     0,
			}
		}

		txIncome := tx.Amount

		// calculate transaction income in local currency
		if tx.Currency != c.business.BaseCurrency {
			rate, err := c.rates.Rate(ctx, tx.Time, tx.Currency, c.business.BaseCurrency)
			if err != nil {
				return fmt.Errorf("get currency rate %s %q -> %q: %w", tx.Time, tx.Currency, c.business.BaseCurrency, err)
			}
			txIncome = Round(txIncome*rate, c.business.RoundingPrecision)
		}

		// accumulate current quarter
		q.Income += txIncome

		return nil
	})
	if err != nil {
		return err
	}

	if q != nil {
		emit()
	}
	return nil
}
