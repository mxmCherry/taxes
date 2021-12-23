package tax

import "fmt"

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
