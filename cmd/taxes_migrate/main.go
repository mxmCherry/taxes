package main

import (
	"math/big"
	"os"
	"time"

	"github.com/mxmCherry/taxes/v2/internal/tax"
	"gopkg.in/yaml.v2"
)

func main() {
	in := yaml.NewDecoder(os.Stdin)
	out := yaml.NewEncoder(os.Stdout)
	defer out.Close()

	var old struct {
		Company struct {
			BaseCurrency string     `yaml:"base_currency"`
			TaxRate      *big.Float `yaml:"tax_rate"`
		} `yaml:"company"`
		Payments []struct {
			Date     time.Time  `yaml:"date"`
			Currency string     `yaml:"currency"`
			Amount   *big.Float `yaml:"amount"`
		} `yaml:"payments"`
	}
	if err := in.Decode(&old); err != nil {
		panic(err.Error())
	}

	run := &tax.CalcRun{
		Calc: tax.Calc{
			LocalCurrency: old.Company.BaseCurrency,
			TaxRate:       old.Company.TaxRate,
		},
	}

	for _, p := range old.Payments {
		tx := &tax.Transaction{
			Time:     p.Date,
			Amount:   p.Amount,
			Currency: p.Currency,
		}

		if len(run.Data) == 0 || run.Data[len(run.Data)-1].Year != tx.Time.Year() {
			run.Data = append(run.Data, &tax.Year{
				Year: tx.Time.Year(),
				Quarters: []*tax.Quarter{{
					Quarter:      QuarterOf(tx.Time),
					Transactions: []*tax.Transaction{tx},
				}},
			})
			continue
		}

		curYear := run.Data[len(run.Data)-1]
		curQuarter := curYear.Quarters[len(curYear.Quarters)-1]
		txQuarter := QuarterOf(tx.Time)
		if txQuarter == curQuarter.Quarter {
			curQuarter.Transactions = append(curQuarter.Transactions, tx)
			continue
		}

		curYear.Quarters = append(curYear.Quarters, &tax.Quarter{
			Quarter:      txQuarter,
			Transactions: []*tax.Transaction{tx},
		})
	}

	if err := out.Encode(run); err != nil {
		panic(err.Error())
	}
}

func QuarterOf(t time.Time) int {
	_, m, _ := t.Date()
	return int(m-1)/3 + 1
}
