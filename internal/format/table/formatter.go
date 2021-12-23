package table

import (
	"fmt"
	"io"
	"math/big"
	"text/tabwriter"

	"github.com/mxmCherry/taxes/v2/internal/format"
	"github.com/mxmCherry/taxes/v2/internal/tax"
)

type formatter struct {
	tab *tabwriter.Writer
}

func New(w io.Writer) format.Formatter {
	return &formatter{
		tab: tabwriter.NewWriter(w, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug),
	}
}

func (f *formatter) Format(c *tax.CalcRun) error {
	_, _ = fmt.Fprintf(f.tab, "Year \t\t Cumulative Income (%s) \t Cumulative Tax (%s)\n",
		c.LocalCurrency,
		c.LocalCurrency,
	)
	_, _ = fmt.Fprintf(f.tab, "\t\t\t\n")

	for _, y := range c.Data {
		for qi, q := range y.Quarters {
			if qi == 0 {
				_, _ = fmt.Fprintf(f.tab, "%d \t Q%d \t % 10s \t % 10s\n",
					y.Year,
					q.Quarter,
					formatFloat(q.CumulativeIncome, c.RoundingPrecision),
					formatFloat(q.CumulativeTax, c.RoundingPrecision),
				)
			} else {
				_, _ = fmt.Fprintf(f.tab, "\t Q%d \t % 10s \t % 10s\n",
					q.Quarter,
					formatFloat(q.CumulativeIncome, c.RoundingPrecision),
					formatFloat(q.CumulativeTax, c.RoundingPrecision),
				)
			}
		}
		_, _ = fmt.Fprintf(f.tab, "\t\t\t\n")
	}

	return nil
}

func (f *formatter) Close() error {
	return f.tab.Flush()
}

func formatFloat(f *big.Float, prec uint) string {
	if prec == 0 {
		return f.String()
	}
	return f.Text('f', int(prec))
}
