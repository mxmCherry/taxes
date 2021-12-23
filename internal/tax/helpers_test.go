package tax_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = DescribeTable("QuarterOf",
	func(m time.Month, q int) {
		Expect(QuarterOf(m)).To(Equal(q))
	},

	Entry("January", time.January, 1),
	Entry("February", time.February, 1),
	Entry("March", time.March, 1),

	Entry("April", time.April, 2),
	Entry("May", time.May, 2),
	Entry("June", time.June, 2),

	Entry("July", time.July, 3),
	Entry("August", time.August, 3),
	Entry("September", time.September, 3),

	Entry("October", time.October, 4),
	Entry("November", time.November, 4),
	Entry("December", time.December, 4),
)

var _ = Describe("Round", func() {
	It("rounds", func() {
		Expect(Round(5.555, 0)).To(Equal(5.555))
		Expect(Round(5.555, 1)).To(Equal(5.6))
		Expect(Round(5.555, 2)).To(Equal(5.56))

		Expect(Round(4.444, 0)).To(Equal(4.444))
		Expect(Round(4.444, 1)).To(Equal(4.4))
		Expect(Round(4.444, 2)).To(Equal(4.44))
	})
})
