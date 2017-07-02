package taxes_test

import (
	"time"

	. "github.com/mxmCherry/taxes/internal/taxes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calc", func() {
	var subject *Calc
	var rates *mockCurrencyRates
	var company Company

	BeforeEach(func() {
		rates = &mockCurrencyRates{
			rate: 26.08,
		}
		subject = NewCalc(rates)

		company = Company{
			BaseCurrency: "UAH",
			TaxRate:      0.05,
		}
	})

	It("should convert payments", func() {
		p := Payment{
			Date:     time.Date(2017, time.March, 1, 0, 0, 0, 0, time.UTC),
			Currency: "USD",
			Amount:   100.00,
		}

		taxes, err := subject.Calc(company, []Payment{p})
		Expect(err).NotTo(HaveOccurred())

		Expect(taxes).To(Equal([]Tax{
			Tax{
				Quarter: Quarter{Year: 2017, Quarter: 1},
				Payments: []BaseCurrencyPayment{
					BaseCurrencyPayment{
						Payment: p,
						Rate:    26.08,
						Amount:  2608.00, // 100.00 * 26.08
					},
				},
				Income: 2608.00,
				Tax:    130.40, // 2608.00 * 5%
			},
		}))

		Expect(rates.request.from).To(Equal("USD"))
		Expect(rates.request.to).To(Equal("UAH"))
		Expect(rates.request.date).To(Equal(p.Date))
	})

	It("should sum quarterly payments", func() {
		p1 := Payment{
			Date:     time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			Currency: "UAH",
			Amount:   100.00,
		}
		p2 := Payment{
			Date:     time.Date(2017, time.March, 1, 0, 0, 0, 0, time.UTC),
			Currency: "USD",
			Amount:   100.00,
		}

		taxes, err := subject.Calc(company, []Payment{p1, p2})
		Expect(err).NotTo(HaveOccurred())

		Expect(taxes).To(Equal([]Tax{
			Tax{
				Quarter: Quarter{Year: 2017, Quarter: 1},
				Payments: []BaseCurrencyPayment{
					BaseCurrencyPayment{
						Payment: p1,
						Rate:    1,
						Amount:  100.00,
					},
					BaseCurrencyPayment{
						Payment: p2,
						Rate:    26.08,
						Amount:  2608.00, // 100.00 * 26.08
					},
				},
				Income: 2708.00, // 100.00 + 2608.00
				Tax:    135.40,  // 2708.00 * 5%
			},
		}))

		Expect(rates.request.from).To(Equal("USD"))
		Expect(rates.request.to).To(Equal("UAH"))
		Expect(rates.request.date).To(Equal(p2.Date))
	})

	It("should calculate taxes for multiple quarters", func() {
		p1 := Payment{
			Date:     time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			Currency: "UAH",
			Amount:   100.00,
		}
		p2 := Payment{
			Date:     time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC),
			Currency: "USD",
			Amount:   100.00,
		}

		taxes, err := subject.Calc(company, []Payment{p1, p2})
		Expect(err).NotTo(HaveOccurred())

		Expect(taxes).To(Equal([]Tax{
			Tax{
				Quarter: Quarter{Year: 2017, Quarter: 1},
				Payments: []BaseCurrencyPayment{
					BaseCurrencyPayment{
						Payment: p1,
						Rate:    1,
						Amount:  100.00,
					},
				},
				Income: 100.00,
				Tax:    5.00, // 100.00 * 5%
			},
			Tax{
				Quarter: Quarter{Year: 2017, Quarter: 2},
				Payments: []BaseCurrencyPayment{
					BaseCurrencyPayment{
						Payment: p2,
						Rate:    26.08,
						Amount:  2608.00, // 100.00 * 26.08
					},
				},
				Income: 2608.00, // 100.00 + 2608.00
				Tax:    130.40,  // 2608.00 * 5%
			},
		}))

		Expect(rates.request.from).To(Equal("USD"))
		Expect(rates.request.to).To(Equal("UAH"))
		Expect(rates.request.date).To(Equal(p2.Date))
	})
})
