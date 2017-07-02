package taxes

// BaseCurrencyPayment holds payment details, converted to company base currency.
type BaseCurrencyPayment struct {
	// Payment holds original payment details.
	Payment Payment
	// Rate holds payment currenty to base currency exchange rate on payment date.
	Rate float64
	// Amount holds payment currency, converted to company base currency.
	Amount float64
}

func CalcBaseCurrencyPayment(rates CurrencyRates, baseCurrency string, p Payment) (*BaseCurrencyPayment, error) {
	if p.Currency == baseCurrency {
		return &BaseCurrencyPayment{
			Payment: p,
			Rate:    1,
			Amount:  p.Amount,
		}, nil
	}

	rate, err := rates.Rate(p.Currency, baseCurrency, p.Date)
	if err != nil {
		return nil, err
	}
	return &BaseCurrencyPayment{
		Payment: p,
		Rate:    rate,
		Amount:  p.Amount * rate,
	}, nil
}
