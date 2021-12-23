# taxes

[![Go Reference](https://pkg.go.dev/badge/github.com/mxmCherry/taxes.svg)](https://pkg.go.dev/github.com/mxmCherry/taxes)
[![Test](https://github.com/mxmCherry/taxes/actions/workflows/test.yml/badge.svg?branch=v2)](https://github.com/mxmCherry/taxes/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mxmCherry/taxes)](https://goreportcard.com/report/github.com/mxmCherry/taxes)

Commandline quarterly tax calculator for simplified Ukrainian taxation system

# Installing

```bash
go install github.com/mxmCherry/taxes/v2@latest
```

# Using

```shell
# basic table output - cumulative income/taxes
taxes taxes.yaml

# advanced YAML output to trace every transaction,
# see Data Format section down the README
taxes -format yaml taxes.yaml
```

## Data Format (`taxes.yaml`)

```yaml
local_currency: UAH   # input: main/local business currency
tax_rate: "0.05"      # input: tax rate, 5%
rounding_precision: 2 # input: round to 2 decimals after comma (to kopecks); do not specify or set to 0 to disable rounding
data:
    - year: 2021               # input: tax year
      total_income: "31438.73" # output (-format yaml): total income amount (in local currency) for this entire year
      total_tax: "1571.95"     # output (-format yaml): total tax amount (in local currency) for this entire year
      quarters:
        - quarter: 1                   # input: quarter index, 1..4
          total_income: "9790.58"      # output (-format yaml): total income amount (in local currency) for this quarter
          total_tax: "489.53"          # output (-format yaml): total tax amount (in local currency) for this quarter
          cumulative_income: "9790.58" # output (-format yaml): cumulative income amount (in local currency) since the beginning of the year
          cumulative_tax: "489.53"     # output (-format yaml): cumulative tax amount (in local currency) since the beginning of the year
          transactions:
            - time: 2021-01-01T00:00:00Z # input: transaction time/date, only date matters
              amount: "100.12"           # input: transaction amount (possibly in foreign currency)
              currency: USD              # input: transaction currency code (possibly foreign currency)
              currency_rate: "28.2746"   # output (-format yaml) (but can be input as well): bank.gov.ua's currency rate for given transaction date/currency
              local_amount: "2830.85"    # output (-format yaml): transaction amount, converted to local currency (basically just `amount` * `currency_rate`)
              tax_rate: "0.05"           # output (-format yaml; can be input as well): tax rate can be overridden per transaction (for example, if business tax rate changed within the year etc)
              tax_amount: "141.54"       # output (-format yaml): tax amount for this transaction (basically just `local_amount` * `tax_rate`)
            ...
        ...
```

Working input example can be found [here](internal/tax/testdata/golden-input-with-rounding.yaml).

By default, this tool returns simplified table output with cumulative income/taxes (since beginning of the year, so Q2 already includes Q1, Q3 includes Q2 and so on) - exactly what's used to declare income:

```
Year |    | Cumulative Income (UAH) | Cumulative Tax (UAH)
     |    |                         |
2019 | Q1 |                    1.00 |       0.05
     | Q2 |                    2.00 |       0.10
     | Q4 |                    3.00 |       0.15
     |    |                         |
2021 | Q2 |              1000000.00 |   50000.00
     | Q4 |             10000000.00 |  500000.00
     |    |                         |
```

You can use `-format yaml` to enable detailed output (as documented above) to trace every transaction detail (currency rates etc):

```shell
taxes -format yaml taxes.yaml
```

## Notes

This CLI tool reads/parses YAML-encoded input file data into memory, pulls/calculates missing data (currency rates, income, tax amounts etc) and prints result to STDOUT.

Transaction-level `currency_rate` can be specified in input file: it makes sense only if number of transactions is quite high and re-pulling all their rates takes too long.
You can get applied currency rates using `-format yaml` output.

v2's CLI API is significantly reduced: no more default locations etc.
Every run requires input file to be explicitly provided as an argument.

It is recommended to comment out previous year data or keep each year in own file to pull currency rates less frequently: that's the slowest bit, and now this tool includes hardcoded 1 RPS rate-limiting to query [bank.gov.ua APIs](https://bank.gov.ua/ua/open-data/api-dev).

## Rounding

Providing non-zero `rounding_precision` (decimals after comma) enables rounding.
It makes sense to use `rounding_precision: 2` to round to kopecks.
Rounding is done for EACH mathematical operation on monetary data.
Basically - after every multiplication (currency conversion) or addition (summing amounts).
Currency rates, provided by [bank.gov.ua APIs](https://bank.gov.ua/ua/open-data/api-dev) are NEVER rounded, they are kept precise (like `33.3333` or so), rounding is done only after foreign currency amount is multiplied by currency rate.

Rounding behavior is the usual one: `x >= 0.5` is rounded up, `x < 0.5` is rounded down:

- given `rounding_precision: 2`
- and transaction `amount: 2.00`
- and transaction `currency_rate: 33.3333`
- result will be `local_amount: 66.67` (and not `66.6666`)

This loss of precision can accumulate across transactions - something to keep in mind.
But that's how most people (and banks) do monetary calculations anyway.

Rounding can be disabled by not specifying it (or specifying `rounding_precision: 0`).

# Migrate from previous major version ([v1](https://github.com/mxmCherry/taxes/tree/v1.0.0))

```shell
go install github.com/mxmCherry/taxes/v2/cmd/taxes_migrate@latest
cat ~/.taxes/data.yaml | taxes_migrate > taxes.yaml
```

Double-check everything after migration as migrate helper is quick && dirty && has no automated tests.

Better keep v1 data file in case anything goes wrong.

Previous version [v1](https://github.com/mxmCherry/taxes/tree/v1.0.0) can still be installed with:

```shell
go install github.com/mxmCherry/taxes@v1.0.0
```
