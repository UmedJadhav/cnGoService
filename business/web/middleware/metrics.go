package middleware

import (
	"cnGoService/business/system/metrics"
	"cnGoService/foundation/web"
	"context"
	"net/http"
)

func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)
			err := handler(ctx, w, r)

			metrics.AddRequests(ctx)
			metrics.AddGoRoutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}
			return err
		}
		return h
	}

	return m
}
