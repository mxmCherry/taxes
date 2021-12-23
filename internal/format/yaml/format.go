package yaml

import (
	"io"

	"github.com/mxmCherry/taxes/v2/internal/format"
	"github.com/mxmCherry/taxes/v2/internal/tax"
	"gopkg.in/yaml.v3"
)

type formatter struct {
	yaml *yaml.Encoder
}

func New(w io.Writer) format.Formatter {
	return &formatter{
		yaml: yaml.NewEncoder(w),
	}
}

func (f *formatter) Format(c *tax.CalcRun) error {
	return f.yaml.Encode(c)
}

func (f *formatter) Close() error {
	return f.yaml.Close()
}
