# taxes

[![Go Reference](https://pkg.go.dev/badge/github.com/mxmCherry/taxes/v2.svg)](https://pkg.go.dev/github.com/mxmCherry/taxes/v2)
[![Test](https://github.com/mxmCherry/taxes/actions/workflows/test.yml/badge.svg)](https://github.com/mxmCherry/taxes/actions/workflows/test.yml)

Commandline quarterly tax calculator for simplified Ukrainian taxation system

# Installing

```shell
go install github.com/mxmCherry/taxes/v2@latest
```

To call the binary without specifying full path, `$GOBIN` can be included in your `$PATH`:

```shell
export GOBIN=$HOME/go/bin # GOBIN defaults to ~/go/bin
export PATH=$PATH:$GOBIN
```

# Using

```shell
# works with ~/.taxes/data.yaml by default
taxes

# or file(s) can be explicitly specified
taxes data1.yaml data2.yaml
```

Input file must be a YAML-encoded data with the following format:

```yaml
business:
  base_currency: UAH    # base/local business currency
  tax_rate: 0.05        # tax rate, 5%
  rounding_precision: 2 # round to 2 decimals after comma (to kopecks); do not specify or set to 0 to disable rounding
transactions:
  - time: 2021-01-01T09:45:53Z # transaction time/date
    amount: 100.12             # transaction amount (possibly in foreign currency)
    currency: USD              # transaction currency code (possibly foreign currency)
  ...
```

Every bit of data mentioned above must be provided.

Transactions must be ordered by time, oldest to newest.

This tool prints (STDOUT) simplified table output with cumulative income/taxes - since beginning of the year, so Q2 already includes Q1, Q3 includes Q2 and so on.

```
          Year |            QQ |        Income |           Tax |
               |               |               |               |
          2020 |            Q1 |        100.00 |          5.00 |
          2020 |            Q3 |        200.00 |         10.00 |
               |               |               |               |
          2021 |            Q1 |        100.12 |          5.01 |
```


## Notes

This CLI tool reads/parses YAML-encoded input file data into memory, pulls/calculates missing data (currency rates, income, tax amounts etc) and prints result to STDOUT.

It is recommended to comment out previous year data or keep each year in own file to pull currency rates less frequently: that's the slowest bit, and now this tool includes hardcoded rate-limiting to query [bank.gov.ua APIs](https://bank.gov.ua/ua/open-data/api-dev) (~150ms between requests).

## Rounding

Providing non-zero `rounding_precision` (decimals after comma) enables rounding.
It makes sense to use `rounding_precision: 2` to round to kopecks.

Rounding is applied in the following cases:

- each currency conversion, so only to non-`base_currency` transactions
- currency rates themselves are not rounded at all
- local-currency transactions are not rounded at all
- rounding is applied to quarter, only in the end (once transaction amounts are summed)
- quarter rounding is applied to income first
- and then it is applied to tax (using already-rounded income amount)

Rounding behavior is the usual one: `x >= 0.5` is rounded up, `x < 0.5` is rounded down.

Rounding can be disabled by not specifying it (or specifying `rounding_precision: 0`).

# Migrate from previous major version ([v1](https://github.com/mxmCherry/taxes/tree/v1.0.0))

Simply rename:

- `company:` -> `business:`
- `payments:` -> `transactions:`
- `date:` -> `time:` (for each transaction)

[v1](https://github.com/mxmCherry/taxes/tree/v1.0.0) can still be installed with:

```shell
go install github.com/mxmCherry/taxes@v1.0.0
```
