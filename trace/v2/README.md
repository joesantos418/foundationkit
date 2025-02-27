# Trace v2

This guide will help you to setup and use the trace/v2 lib. We have a usage example in `examples` folder.

## Configuring

This lib provides a `Config` struct and this section will explain how to properly configure.

### Provider

1. `stackdriver` to export the trace to stackdriver
2. Empty string to not export the trace, but the trace will be created

### Probability sample

1. `0` will not sample or export (if provider was set)
2. Float between 0 and 1 will use the number as a probability for sampling and export (if provider was set)
3. `1` will sample and export (if provider was set)

### Stackdriver

Information to connect in stackdriver

## Setting up and using the lib

Using in your config.go file

```golang
var config struct {
    // (...)

	Trace trace.Config
}
```

Creating a setup method

```golang
func setupTrace() {
	traceShutdown := trace.Setup(config.Trace)
	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "opentelemetry_trace",
			Priority: app.ShutdownPriority(shutdownPriorityTrace),
			Handler:  traceShutdown,
			Policy:   app.ErrorPolicyAbort,
		})
}
```

Setting the trace up AFTER create a new App in main.go

```golang

func main() {
	app.SetupConfig(&config)

    // (...)

	if err := app.NewDefaultApp(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to create app")
	}

	setupTrace()

    // (...)
}
```

Starting a span

```golang
ctx, span := trace.Start(ctx, "SPAN-NAME")
defer span.End()
```

Recovering trace informations and logging it

```golang
type TraceInfo struct {
	ID        string
	IsSampled bool
}
```

```golang
t := trace.GetTraceInfoFromContext(ctx)
log.Ctx(ctx).Info().EmbedObject(t).Msg("Hello")
```

## Propagating the trace

### API

Using in transport layer

```golang
func MakeHTTPHandler(e endpoint.Endpoint) http.Handler {
    // (...)

	r := mux.NewRouter()
	r.Use(trace.MuxHTTPMiddleware("SERVER-NAME"))

    //(...)
	return r
}
```

Using in a HTTP request

```golang
request, err := http.NewRequestWithContext(
    ctx,
    "POST",
    url,
    bytes.NewReader(body),
)
if err != nil {
    return "", errors.E(op, err)
}

trace.SetTraceInRequest(request)
```

Using in a HTTP response

```golang
func encodeResponse(
	ctx context.Context,
	w http.ResponseWriter,
	r interface{},
) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	trace.SetTraceInResponse(ctx, w)
	return json.NewEncoder(w).Encode(r)
}
```

### Workers

[WIP]