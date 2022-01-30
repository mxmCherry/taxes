package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/mxmCherry/taxes/v2/internal/tax"
)

type table struct {
	tab           *tabwriter.Writer
	headerWritten bool
	prevYear      int
}

func Table(w io.Writer) Formatter {
	return &table{
		tab: tabwriter.NewWriter(w, 15, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug),
	}
}

func (f *table) Format(biz *tax.Business, q *tax.Quarter) error {
	if !f.headerWritten {
		if _, err := fmt.Fprintf(f.tab, "Year \tQQ \tIncome \tTax \tCum Income \tCum Tax \t\n"); err != nil {
			return err
		}
		f.headerWritten = true
	}

	if q.Year != f.prevYear {
		if _, err := fmt.Fprintf(f.tab, "\t\t\t\t\t\t\n"); err != nil {
			return err
		}
		f.prevYear = q.Year
	}

	var err error
	if biz.RoundingPrecision > 0 {
		_, err = fmt.Fprintf(f.tab, "%d \tQ%d \t% .*f \t% .*f \t% .*f \t% .*f \t\n",
			q.Year,
			q.Quarter,
			biz.RoundingPrecision,
			q.Income,
			biz.RoundingPrecision,
			q.Tax,
			biz.RoundingPrecision,
			q.CumIncome,
			biz.RoundingPrecision,
			q.CumTax,
		)
	} else {
		_, err = fmt.Fprintf(f.tab, "%d \tQ%d \t%f \t%f \t%f \t%f \t\n",
			q.Year,
			q.Quarter,
			q.Income,
			q.Tax,
			q.CumIncome,
			q.CumTax,
		)
	}
	return err
}

func (f *table) Close() error {
	return f.tab.Flush()
}
