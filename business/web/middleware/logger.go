package middleware

import (
	"cnGoService/foundation/web"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)
			if err != nil {
				return err //web.NewShutdownError("web values missing context")
			}
			log.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)
			err = handler(ctx, w, r)
			log.Infow("request completed", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))
			return err
		}
		return h
	}
	return m
}
