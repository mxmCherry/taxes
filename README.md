# taxes

[![Go Reference](https://pkg.go.dev/badge/github.com/mxmCherry/taxes.svg)](https://pkg.go.dev/github.com/mxmCherry/taxes)
[![Test](https://github.com/mxmCherry/taxes/actions/workflows/test.yml/badge.svg?branch=v2)](https://github.com/mxmCherry/taxes/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mxmCherry/taxes)](https://goreportcard.com/report/github.com/mxmCherry/taxes)

Commandline quarterly tax calculator for simplified Ukrainian taxation system

# Install

```bash
go install github.com/mxmCherry/taxes@v2
```

# Use

```shell
taxes taxes.yaml
```

For input file format (`taxes.yaml` file in the snippet above), refer to [test input example](internal/tax/testdata/golden-input-with-rounding.yaml).

For output format, refer to [test output example](internal/tax/testdata/golden-output-with-rounding.yaml) (pretty much the same as input, but with additional fields pulled/calculated).

# Notes

This CLI tool reads/parses YAML-encoded input file data into memory, pulls/calculates missing data (currency rates, income, tax amounts etc) and prints YAML result to STDOUT.

Transaction-level `currency_rate` can be specified in input: it makes sense only if number of transactions is quite high and re-pulling all their rates takes too long.
You can get used currency rates from output (same format as input, but more fields are populated).

v2's CLI API is significantly reduced: no more default locations etc.
Every run requires input file to be explicitly provided as an argument.

# Rounding

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
go install github.com/mxmCherry/taxes/cmd/taxes_migrate@v2
cat ~/.taxes/data.yaml | taxes_migrate > ~/.taxes/data-v2.yaml
```

Double-check everything after migration as migrate helper is quick && dirty && has no automated tests.
