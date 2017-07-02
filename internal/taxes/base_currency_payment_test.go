package taxes_test

import (
	"time"

	. "github.com/mxmCherry/taxes/internal/taxes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalcBaseCurrencyPayment", func() {
	var rates *mockCurrencyRates

	BeforeEach(func() {
		rates = &mockCurrencyRates{
			rate: 26.08,
		}
	})

	It("should handle payments in base currency", func() {
		payment := Payment{
			Date:     time.Date(2017, time.July, 2, 14, 40, 0, 0, time.UTC),
			Currency: "UAH",
			Amount:   100.00,
		}

		actual, err := CalcBaseCurrencyPayment(rates, "UAH", payment)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).NotTo(BeNil())

		Expect(actual.Payment).To(Equal(payment))
		Expect(actual.Rate).To(Equal(float64(1)))
		Expect(actual.Amount).To(Equal(float64(100.00)))
	})

	It("should convert payments in foreign currency", func() {
		payment := Payment{
			Date:     time.Date(2017, time.July, 2, 14, 40, 0, 0, time.UTC),
			Currency: "USD",
			Amount:   100.00,
		}

		actual, err := CalcBaseCurrencyPayment(rates, "UAH", payment)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).NotTo(BeNil())

		Expect(actual.Payment).To(Equal(payment))
		Expect(actual.Rate).To(Equal(float64(26.08)))
		Expect(actual.Amount).To(Equal(float64(2608.00)))

		Expect(rates.request.from).To(Equal("USD"))
		Expect(rates.request.to).To(Equal("UAH"))
		Expect(rates.request.date).To(Equal(payment.Date))
	})

})
