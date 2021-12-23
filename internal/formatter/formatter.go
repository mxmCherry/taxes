package formatter

import "github.com/mxmCherry/taxes/v2/internal/tax"

type Formatter interface {
	Format(*tax.Business, *tax.Quarter) error
	Close() error
}
