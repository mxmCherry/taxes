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

TODO: finish it and rewrite doc

# Migrate from previous major version ([v1](https://github.com/mxmCherry/taxes/tree/v1.0.0))

```shell
go install github.com/mxmCherry/taxes/cmd/taxes_migrate@v2
cat ~/.taxes/data.yaml | taxes_migrate > taxes-v2.yaml
```

Double-check everything after migration as migrate helper is quick && dirty && has no automated tests.
