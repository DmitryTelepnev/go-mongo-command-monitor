# MongoDB Command monitor

[![Go Report Card](https://goreportcard.com/badge/github.com/DmitryTelepnev/mongo-command-monitor)](https://goreportcard.com/report/github.com/DmitryTelepnev/mongo-command-monitor)

Lib for GoLang MongoDB driver for command monitoring

Implements:
* total number of started queries with "app", "database", "collection", "command" labels
* total number of finished queries with "app", "database", "collection", "command" labels
* mongo query duration with "app", "database", "collection", "command" labels

Also, you can implement a custom metric collector by the Metrics interface.

## Import into your project

```bash
go get github.com/DmitryTelepnev/mongo-command-monitor
```

```
import "github.com/DmitryTelepnev/mongo-command-monitor"

opts := driverOpts.Client()
...
opts.SetMonitor(commandMonitor.GetCommandMonitor(metrics.NewPrometheus("examples-app")))
...
client, err := driver.Connect(connectCtx, opts)
```

## Run examples

```bash
cd examples; docker-compose up --build --force-recreate
```

## Run benchmarks

```bash
cd benchmarks; docker-compose up --build --force-recreate --remove-orphans -V --abort-on-container-exit
```
## Benchmark results

```text
benchmark_1                 | BenchmarkMongoInsertQueries_WithCommandMonitor-4              2767            387822 ns/op            5554 B/op        103 allocs/op
benchmark_1                 | BenchmarkMongoInsertQueries_WithoutCommandMonitor-4           2982            353121 ns/op            4388 B/op         75 allocs/op
benchmark_1                 | BenchmarkMongoUpdateQueries_WithCommandMonitor-4              3511            381561 ns/op            7279 B/op        130 allocs/op
benchmark_1                 | BenchmarkMongoUpdateQueries_WithoutCommandMonitor-4           3253            337433 ns/op            5797 B/op        101 allocs/op
benchmark_1                 | BenchmarkMongoFindQueries_WithCommandMonitor-4                4606            280433 ns/op            5999 B/op         98 allocs/op
benchmark_1                 | BenchmarkMongoFindQueries_WithoutCommandMonitor-4             4932            241254 ns/op            4947 B/op         71 allocs/op
```
