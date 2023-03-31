package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handler) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)
	statusCode := http.StatusOK
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)
}
