package testgrp

import (
	"cnGoService/foundation/web"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handler) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)
	return web.Respond(ctx, w, status, statusCode)
}
