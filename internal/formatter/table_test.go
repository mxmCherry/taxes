package formatter_test

import (
	"bytes"

	"github.com/mxmCherry/taxes/v2/internal/tax"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/formatter"
)

var _ = Describe("Table", func() {
	It("formats", func() {
		biz := &tax.Business{RoundingPrecision: 2}

		buf := bytes.NewBuffer(nil)
		subject := Table(buf)

		Expect(subject.Format(biz, &tax.Quarter{Year: 2020, Quarter: 1, Income: 100.00, Tax: 5.00})).To(Succeed())
		Expect(subject.Format(biz, &tax.Quarter{Year: 2020, Quarter: 3, Income: 200.00, Tax: 10.00})).To(Succeed())
		Expect(subject.Format(biz, &tax.Quarter{Year: 2021, Quarter: 1, Income: 10.00, Tax: 0.50})).To(Succeed())
		Expect(subject.Close()).To(Succeed())

		Expect(buf.String()).To(Equal("" +
			"          Year |            QQ |        Income |           Tax |\n" +
			"               |               |               |               |\n" +
			"          2020 |            Q1 |        100.00 |          5.00 |\n" +
			"          2020 |            Q3 |        200.00 |         10.00 |\n" +
			"               |               |               |               |\n" +
			"          2021 |            Q1 |         10.00 |          0.50 |\n" +
			"",
		))
	})
})
