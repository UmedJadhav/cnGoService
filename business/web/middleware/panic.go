package middleware

import (
	"cnGoService/business/system/metrics"
	"cnGoService/foundation/web"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
)

func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// defer recovers from panic and set the err return variable after panic
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, trace)
					metrics.AddPanics(ctx)
				}
			}()

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
