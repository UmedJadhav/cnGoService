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
	return web.Respond(ctx, w, status, statusCode)
}
