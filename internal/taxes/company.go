package taxes

// Company holds company details.
type Company struct {
	// BaseCurrency holds company base (tax) currency, ISO 4217 (e.g., UAH).
	BaseCurrency string `yaml:"base_currency"`
	// TaxRate holds company tax rate (e.g., %5 - 0.05).
	TaxRate float64 `yaml:"tax_rate"`
}
