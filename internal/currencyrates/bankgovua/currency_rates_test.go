package bankgovua_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/mxmCherry/taxes/internal/currencyrates/bankgovua"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CurrencyRates", func() {
	var subject *CurrencyRates
	var server *httptest.Server
	var lastRequest struct {
		from string
		date string
	}

	BeforeEach(func() {
		lastRequest.from = ""
		lastRequest.date = ""

		mux := http.NewServeMux()
		mux.HandleFunc("/NBUStatService/v1/statdirectory/exchange", func(w http.ResponseWriter, r *http.Request) {
			lastRequest.from = r.URL.Query().Get("valcode")
			lastRequest.date = r.URL.Query().Get("date")
			_, _ = io.WriteString(w, `[{"rate":26.08}]`)
		})
		server = httptest.NewServer(mux)

		subject = NewCurrencyRatesWithAddr(server.URL)
	})

	AfterEach(func() {
		server.Close()
	})

	It("should return rates", func() {
		rate, err := subject.Rate("USD", "UAH", time.Date(2017, time.July, 2, 15, 35, 0, 0, time.UTC))
		Expect(err).NotTo(HaveOccurred())
		Expect(rate).To(Equal(float64(26.08)))

		Expect(lastRequest.from).To(Equal("USD"))
		Expect(lastRequest.date).To(Equal("20170702"))
	})

	It("should cache rates", func() {
		rate1, err := subject.Rate("USD", "UAH", time.Date(2017, time.July, 2, 15, 35, 0, 0, time.UTC))
		Expect(err).NotTo(HaveOccurred())
		Expect(rate1).To(Equal(float64(26.08)))

		Expect(lastRequest.from).NotTo(BeEmpty())
		Expect(lastRequest.date).NotTo(BeEmpty())

		lastRequest.from = ""
		lastRequest.date = ""

		rate2, err := subject.Rate("USD", "UAH", time.Date(2017, time.July, 2, 19, 48, 0, 0, time.UTC))
		Expect(err).NotTo(HaveOccurred())
		Expect(rate2).To(Equal(float64(26.08)))

		Expect(lastRequest.from).To(BeEmpty())
		Expect(lastRequest.date).To(BeEmpty())

		rate3, err := subject.Rate("USD", "UAH", time.Date(2017, time.June, 2, 19, 48, 0, 0, time.UTC))
		Expect(err).NotTo(HaveOccurred())
		Expect(rate3).To(Equal(float64(26.08)))

		Expect(lastRequest.from).NotTo(BeEmpty())
		Expect(lastRequest.date).NotTo(BeEmpty())
	})

})
