package tax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = Describe("Transaction", func() {
	DescribeTable("Validate",
		func(tx *Transaction, errStr string) {
			if errStr == "" {
				Expect(tx.Validate()).To(Succeed())
			} else {
				Expect(tx.Validate()).To(MatchError(errStr))
			}
		},

		Entry("nil",
			nil,
			"cannot be nil",
		),
		Entry("bad amount",
			&Transaction{Amount: 0.00, Currency: "UAH"},
			"amount 0.000000 if invalid: should be greater than 0",
		),
		Entry("bad currency",
			&Transaction{Amount: 0.01, Currency: "foo"},
			`currency "foo" is invalid: should be a 3-letter ISO code, uppercased (like "UAH")`,
		),
	)
})
