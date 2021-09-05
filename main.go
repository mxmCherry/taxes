package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mxmCherry/taxes/internal/bankgovua"
	"github.com/mxmCherry/taxes/internal/tax"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"
)

func main() {
	output := yaml.NewEncoder(os.Stdout)
	defer output.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	for _, filename := range os.Args[1:] {
		if err := process(ctx, output, filename); err != nil {
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

func process(ctx context.Context, output *yaml.Encoder, filename string) error {
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

	if err := output.Encode(run); err != nil {
		return fmt.Errorf("output yaml: %w", err)
	}
	return nil
}
