package bankgovua_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/mxmCherry/taxes/v2/internal/bankgovua"
)

var _ = Describe("ParseResponse", func() {
	It("parses valid responses", func() {
		raw := []byte(`[
			{
				"r030": 826,
				"txt":"Фунт стерлінгів",
				"rate":36.2722,
				"cc":"GBP",
				"exchangedate":"23.12.2021"
			}
		]`)
		Expect(ParseResponse(raw)).To(Equal(36.2722))
	})

	DescribeTable("response validation",
		func(raw, errStr string) {
			_, err := ParseResponse([]byte(raw))
			Expect(err).To(MatchError(HavePrefix(errStr)))
		},
		Entry("invalid JSON",
			`INVALID`,
			"decode JSON: ",
		),
		Entry("no rate returned",
			`[]`,
			"unexpected response size: ",
		),
		Entry("too many rates returned",
			`[{}, {}]`,
			"unexpected response size: ",
		),
		Entry("bad rate returned",
			`[{"rate": -1.23}]`,
			"no rate provided",
		),
		Entry("zero rate returned",
			`[{"rate": 0}]`,
			"no rate provided",
		),
	)
})

var _ = Describe("BuildURL", func() {
	It("builds API URL", func() {
		Expect(
			BuildURL(time.Date(2021, time.December, 23, 9, 45, 17, 123456, time.UTC), "GBP"),
		).To(Equal("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=GBP&date=20211223&json"))
	})
})
