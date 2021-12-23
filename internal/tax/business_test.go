package tax_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = Describe("Business", func() {

	DescribeTable("Validate",
		func(biz *Business, errStr string) {
			if errStr == "" {
				Expect(biz.Validate()).To(Succeed())
			} else {
				Expect(biz.Validate()).To(MatchError(errStr))
			}
		},

		Entry("nil",
			nil,
			"cannot be nil",
		),
		Entry("bad currency",
			&Business{BaseCurrency: "foo"},
			`base_currency "foo" is invalid: should be a 3-letter ISO code, uppercased (like "UAH")`,
		),
		Entry("bad tax rate (lower bound)",
			&Business{BaseCurrency: "UAH", TaxRate: 0},
			"tax_rate 0.000000 is invalid: should be in 0..1 range",
		),
		Entry("bad tax rate (upper bound)",
			&Business{BaseCurrency: "UAH", TaxRate: 1},
			"tax_rate 1.000000 is invalid: should be in 0..1 range",
		),
		Entry("valid",
			&Business{BaseCurrency: "UAH", TaxRate: 0.05},
			"",
		),
	)
})
