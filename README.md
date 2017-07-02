# taxes

Commandline quarterly tax calculator for simplified Ukrainian taxation system

# Install

```bash
go get -u github.com/mxmCherry/taxes/cmd/taxes
```

# Use

By default (if no `--file` provided), `taxes` command tries to consume `~/.taxes/data.yaml`.
Refer to [example.yaml](example.yaml) for example structure.

```bash
$ taxes --file=example.yaml
Year	QQ	Payments	Income (UAH)	   Tax (UAH)
2016	Q3	       2	      268.20	       13.41
2016	Q4	       2	      805.99	       40.30
2017	Q1	       2	     1434.60	       71.73
2017	Q2	       2	     1917.27	       95.87
```
