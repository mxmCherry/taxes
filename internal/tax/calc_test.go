package tax_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"gopkg.in/yaml.v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mxmCherry/taxes/internal/tax"
)

var _ = Describe("CalcRun", func() {
	var subject *tax.CalcRun

	BeforeEach(func() {
		subject = new(tax.CalcRun)
		unmarshalYAML("testdata/golden-input.yaml", subject)

		subject.Calc.CurrencyRates = currencyRates{
			"2021-01-01": {
				"USD": {"UAH": "28.2746"},
				"EUR": {"UAH": "34.7396"},
			},
			"2021-07-01": {
				"USD": {"UAH": "27.2275"},
				"EUR": {"UAH": "32.3463"},
			},
		}
	})

	It("works", func() {
		Expect(subject.Run(context.Background())).To(Succeed())

		b, err := yaml.Marshal(subject)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(b)).To(Equal(string(readFile("testdata/golden-output.yaml"))))
	})
})

func unmarshalYAML(filename string, run *tax.CalcRun) {
	b, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())

	Expect(yaml.Unmarshal(b, run)).To(Succeed())
}

func readFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	return b
}

type currencyRates map[string]map[string]map[string]string // currencyRates[YYYY-MM-DD][FROM][TO] = rate

func (rr currencyRates) Rate(_ context.Context, date time.Time, from, to string) (*big.Float, error) {
	dateKey := date.Format("2006-01-02")
	rate := rr[dateKey][from][to]
	if rate == "" {
		return nil, fmt.Errorf("no rate for %q -> %q (%q)", from, to, dateKey)
	}

	r, _, err := big.ParseFloat(rate, 10, 0, big.ToNearestEven)
	Expect(err).NotTo(HaveOccurred())
	return r, nil
}
