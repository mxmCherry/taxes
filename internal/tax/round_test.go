package tax_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = Describe("Round", func() {
	It("rounds", func() {
		Expect(tax.Round(f("1.234"), 2).String()).To(Equal("1.23"))
		Expect(tax.Round(f("1.235"), 2).String()).To(Equal("1.24"))
	})
})

func f(s string) *big.Float {
	x, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return x
}
