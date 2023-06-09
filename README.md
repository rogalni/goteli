[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/opentelemetry-go-extra/otelzap)](https://pkg.go.dev/github.com/uptrace/opentelemetry-go-extra/otelzap)

# GO OpenTelemetry instrumentation  


## Installation

```shell
go get github.com/rogalni/goteli
```

## Usage

You need to create an `goteli.New` with a context and options.

```go
package main

import (
  "context"
  "github.com/rogalni/goteli"
  "github.com/uptrace/opentelemetry-go-extra/otelzap"
  "go.opentelemetry.io/otel"
)

func main() {
  ctx := context.Background()
  opts := goteli.NewDefaultOpts()
  cu := goteli.New(ctx, opts)
  defer cu(ctx)
}

```

