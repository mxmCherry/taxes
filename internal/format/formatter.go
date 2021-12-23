package format

import "github.com/mxmCherry/taxes/v2/internal/tax"

type Formatter interface {
	Format(*tax.CalcRun) error
	Close() error
}
