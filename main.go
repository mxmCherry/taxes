package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mxmCherry/taxes/v2/internal/bankgovua"
	"github.com/mxmCherry/taxes/v2/internal/format"
	"github.com/mxmCherry/taxes/v2/internal/tax"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"

	ftable "github.com/mxmCherry/taxes/v2/internal/format/table"
	fyaml "github.com/mxmCherry/taxes/v2/internal/format/yaml"
)

var formats map[string]func(io.Writer) format.Formatter

var flags struct {
	Format string
}

func init() {
	formats = map[string]func(io.Writer) format.Formatter{
		"yaml":  fyaml.New,
		"table": ftable.New,
	}
	flag.StringVar(&flags.Format, "format", "table", "output format: table|yaml")
}

func main() {
	flag.Parse()

	newFormatter, ok := formats[flags.Format]
	if !ok {
		fmt.Fprintf(os.Stderr, "ERROR: unsupported -format %q", flags.Format)
		os.Exit(1)
	}

	formatter := newFormatter(os.Stdout)
	defer formatter.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	for _, filename := range flag.Args() {
		if err := process(ctx, formatter, filename); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %q: %s\n", filename, err)
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func process(ctx context.Context, formatter format.Formatter, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close()

	run := tax.CalcRun{
		Calc: tax.Calc{
			CurrencyRates: &bankgovua.Client{
				HTTP:        new(http.Client),
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
				MaxRetries:  3,
			},
		},
	}
	if err := yaml.NewDecoder(f).Decode(&run); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}

	if err := run.Run(ctx); err != nil {
		return fmt.Errorf("tax calc: %w", err)
	}

	if err := formatter.Format(&run); err != nil {
		return fmt.Errorf("output yaml: %w", err)
	}
	return nil
}
