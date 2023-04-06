package metrics

import (
	"context"
	"expvar"
)

type metrics struct {
	goRoutines *expvar.Int
	requests   *expvar.Int
	Errors     *expvar.Int
	Panics     *expvar.Int
}

// This holds a single instance of the metrics value needed for collecting metrics.
// The expvar package already has a singleton for different metrics that are registered with the
// package so there isn't much choice there.
var m *metrics

func init() {
	m = &metrics{
		goRoutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		Errors:     expvar.NewInt("errors"),
		Panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

func AddGoRoutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.goRoutines.Value()%100 == 0 {
			v.goRoutines.Add(1)
		}
	}
}

func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)

	}
}

func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.Errors.Add(1)

	}
}

func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.Panics.Add(1)
	}
}
