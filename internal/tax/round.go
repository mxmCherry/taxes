package tax

import "math/big"

// Round rounds provided float to `prec` number of decimals after comma.
func Round(x *big.Float, prec uint) *big.Float {
	// heavy, dirty, but does rounding job just fine
	str := x.Text('f', int(prec)) // already rounded
	res, _, err := big.ParseFloat(str, 10, 0, big.ToNearestEven)
	if err != nil {
		panic("rounding error: " + err.Error()) // impossible
	}
	return res
}
