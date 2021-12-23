package tax

import (
	"fmt"
	"time"
)

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

// ----------------------------------------------------------------------------

type TransactionSlice []*Transaction

func (tt TransactionSlice) Each(cb func(*Transaction) error) error {
	for _, t := range tt {
		if err := cb(t); err != nil {
			return err
		}
	}
	return nil
}
