/*

Command taxes calculates quarterly taxes from company/payments YAML file.

Usage:

	taxes --file=example.yaml

Company/payments file example:

  company:
    base_currency: UAH
    tax_rate:      0.05
  payments:
    - date:     2016-08-09T16:54:00+03:00
      currency: USD
      amount:   10.00
    - date:     2016-08-10T10:33:00+03:00
      currency: UAH
      amount:   20.00

*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/mxmCherry/taxes/internal/currencyrates/bankgovua"
	"github.com/mxmCherry/taxes/internal/taxes"
)

var flags struct {
	file string
}

func init() {
	defaultFilePath := filepath.Join(os.Getenv("HOME"), ".taxes/data.yaml")
	flag.StringVar(&flags.file, "file", defaultFilePath, "Path to data YAML file")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run() error {
	dataBytes, err := ioutil.ReadFile(flags.file)
	if err != nil {
		return err
	}

	var data struct {
		Company  taxes.Company
		Payments []taxes.Payment
	}
	if err := yaml.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	rates := bankgovua.NewCurrencyRates()
	calc := taxes.NewCalc(rates)

	taxes, err := calc.Calc(data.Company, data.Payments)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Year\tQQ\tPayments\tIncome (%s)\t   Tax (%s)\n",
		data.Company.BaseCurrency,
		data.Company.BaseCurrency,
	)
	for _, tax := range taxes {
		fmt.Printf(
			"%4d\tQ%1d\t% 8d\t%12.2f\t%12.2f\n",
			tax.Quarter.Year,
			tax.Quarter.Quarter,
			len(tax.Payments),
			ceil(tax.Income),
			ceil(tax.Tax),
		)
	}
	return nil
}

func ceil(x float64) float64 {
	return math.Ceil(x*100) / 100
}
