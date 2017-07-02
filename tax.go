package taxes

// Tax holds tax details.
type Tax struct {
	// Quarters holds financial year quarter.
	Quarter Quarter `yaml:"quarter"`
	// Payments holds quarterly company bank payments.
	Payments []BaseCurrencyPayment `yaml:"payments"`
	// Income holds quarterly company income.
	Income float64 `yaml:"income"`
	// Tax holds quarterly tax amount.
	Tax float64 `yaml:"tax"`
}

// ----------------------------------------------------------------------------

type taxesByQuarter []Tax

func (t taxesByQuarter) Len() int {
	return len(t)
}

func (t taxesByQuarter) Less(i, j int) bool {
	l, r := t[i], t[j]
	return l.Quarter.Year < r.Quarter.Year ||
		(l.Quarter.Year == r.Quarter.Year && l.Quarter.Quarter < r.Quarter.Quarter)
}

func (t taxesByQuarter) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
