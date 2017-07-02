package taxes_test

import (
	"time"

	. "github.com/mxmCherry/taxes"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("QuarterOf",
	func(t time.Time, q Quarter) {
		Expect(QuarterOf(t)).To(Equal(q))
	},
	Entry("January",
		time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 1},
	),
	Entry("February",
		time.Date(2017, time.February, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 1},
	),
	Entry("March",
		time.Date(2017, time.March, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 1},
	),
	Entry("April",
		time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 2},
	),
	Entry("May",
		time.Date(2017, time.May, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 2},
	),
	Entry("June",
		time.Date(2017, time.June, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 2},
	),
	Entry("July",
		time.Date(2017, time.July, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 3},
	),
	Entry("August",
		time.Date(2017, time.August, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 3},
	),
	Entry("September",
		time.Date(2017, time.September, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 3},
	),
	Entry("October",
		time.Date(2017, time.October, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 4},
	),
	Entry("November",
		time.Date(2017, time.November, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 4},
	),
	Entry("December",
		time.Date(2017, time.December, 1, 0, 0, 0, 0, time.UTC),
		Quarter{Year: 2017, Quarter: 4},
	),
)
