package taxes

// Tax holds tax details.
type Tax struct {
	// Quarters holds financial year quarter.
	Quarter Quarter
	// Payments holds quarterly company bank payments.
	Payments []BaseCurrencyPayment
	// Income holds quarterly company income.
	Income float64
	// Tax holds quarterly tax amount.
	Tax float64
}
