# taxes

[![Go Reference](https://pkg.go.dev/badge/github.com/mxmCherry/taxes.svg)](https://pkg.go.dev/github.com/mxmCherry/taxes)
[![Test](https://github.com/mxmCherry/taxes/actions/workflows/test.yml/badge.svg?branch=v2)](https://github.com/mxmCherry/taxes/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mxmCherry/taxes)](https://goreportcard.com/report/github.com/mxmCherry/taxes)

Commandline quarterly tax calculator for simplified Ukrainian taxation system

# Install

```bash
go install github.com/mxmCherry/taxes/v2@latest
```

# Use

```shell
taxes taxes.yaml
```

For input file format (`taxes.yaml` file in the snippet above), refer to [test input example](internal/tax/testdata/golden-input-with-rounding.yaml).

For output format, refer to [test output example](internal/tax/testdata/golden-output-with-rounding.yaml) (pretty much the same as input, but with additional fields pulled/calculated).

Here's the data format with fields explained:

```yaml
# taxes.yaml:
local_currency: UAH   # input: main/local business currency
tax_rate: "0.05"      # input: tax rate, 5%
rounding_precision: 2 # input: round to 2 decimals after comma (to kopecks); do not specify or set to 0 to disable rounding
data:
    - year: 2021               # input: tax year
      total_income: "31438.73" # output: total income amount (in local currency) for this entire year
      total_tax:     "1571.95" # output: total tax amount (in local currency) for this entire year
      quarters:
        - quarter: 1                   # input: quarter index, 1..4
          total_income: "9790.58"      # output: total income amount (in local currency) for this quarter
          total_tax: "489.53"          # output: total tax amount (in local currency) for this quarter
          cumulative_income: "9790.58" # output: cumulative income amount (in local currency) since the beginning of the year
          cumulative_tax: "489.53"     # output: cumulative tax amount (in local currency) since the beginning of the year
          transactions:
            - time: 2021-01-01T00:00:00Z # input: transaction time/date, only date matters
              amount: "100.12"           # input: transaction amount (possibly in foreign currency)
              currency: USD              # input: transaction currency code (possibly foreign currency)
              currency_rate: "28.2746"   # output (but can be input as well): bank.gov.ua's currency rate for given transaction date/currency
              local_amount: "2830.85"    # output: transaction amount, converted to local currency (basically just `amount` * `currency_rate`)
              tax_rate: "0.05"           # output (but can be input as well): tax rate can be overridden per transaction (for example, if business tax rate changed within the year etc)
              tax_amount: "141.54"       # output: tax amount for this transaction (basically just `local_amount` * `tax_rate`)
            ...
        ...
```

## Notes

This CLI tool reads/parses YAML-encoded input file data into memory, pulls/calculates missing data (currency rates, income, tax amounts etc) and prints YAML result to STDOUT.

Transaction-level `currency_rate` can be specified in input: it makes sense only if number of transactions is quite high and re-pulling all their rates takes too long.
You can get used currency rates from output (same format as input, but more fields are populated).

v2's CLI API is significantly reduced: no more default locations etc.
Every run requires input file to be explicitly provided as an argument.

It is recommended to comment out previous year data or keep each year in own file to pull currency rates less frequently: that's the slowest bit, and now this tool includes hardcoded 1 RPS rate-limiting to query [bank.gov.ua APIs](https://bank.gov.ua/ua/open-data/api-dev).

## Rounding

Usage input example includes rounding to 2 decimals after comma (effectively - round to kopecks/cents/...).
Rounding is done for EACH mathematical operation on monetary data.
Basically - after every multiplication (currency conversion) or addition (summing amounts).

Rounding behaviour is the usual one: `x >= 0.5` is rounded up, `x < 0.5` is rounded down:

- given `rounding_precision: 2`
- and transaction `amount: 2.00`
- and transaction `currency_rate: 33.3333`
- result will be `66.67` (and not `66.6666`)

This loss of precision can accumulate across transactions - something to keep in mind.
But that's how most people (and banks) do monetary calculations anyway.

Rounding can be disabled by not specifying it (or specifying `rounding_precision: 0`).

# Migrate from previous major version ([v1](https://github.com/mxmCherry/taxes/tree/v1.0.0))

```shell
go install github.com/mxmCherry/taxes/v2/cmd/taxes_migrate@latest
cat ~/.taxes/data.yaml | taxes_migrate > taxes.yaml
```

Double-check everything after migration as migrate helper is quick && dirty && has no automated tests.
