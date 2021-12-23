package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/mxmCherry/taxes/v2/internal/bankgovua"
	"github.com/mxmCherry/taxes/v2/internal/formatter"
	"github.com/mxmCherry/taxes/v2/internal/tax"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	out := formatter.Table(os.Stdout)
	defer out.Close()

	rates := &bankgovua.Client{
		HTTP:        new(http.Client),
		RateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
		MaxRetries:  3,
	}

	filenames := flag.Args()
	if len(filenames) == 0 {
		filenames = []string{
			filepath.Join(os.Getenv("HOME"), ".taxes/data.yaml"),
		}
	}

	for _, filename := range filenames {
		if err := process(ctx, out, rates, filename); err != nil {
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

func process(ctx context.Context, out formatter.Formatter, rates tax.CurrencyRates, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close()

	var data struct {
		Business     *tax.Business        `yaml:"business"`
		Transactions tax.TransactionSlice `yaml:"transactions"`
	}

	if err := yaml.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}

	calc, err := tax.NewCalc(
		data.Business,
		data.Transactions,
		rates,
	)
	if err != nil {
		return err
	}

	return calc.Each(ctx, func(q *tax.Quarter) error {
		return out.Format(data.Business, q)
	})
}
