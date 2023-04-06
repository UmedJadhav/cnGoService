package handlers

import (
	"cnGoService/app/services/sales-api/handlers/debug/checkgrp"
	"cnGoService/app/services/sales-api/handlers/v1/testgrp"
	"cnGoService/business/web/middleware"
	"cnGoService/foundation/web"
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from std library into a new mux
// bypassing the use of DefaultServerMux. Using DefaultServerMux would be a security risk since a dependency
// could inject a handler into our service without us knowing about it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)
	return mux
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs an http.Handler with all application routes defined
func APIMux(conf APIMuxConfig) *web.App {
	app := web.NewApp(
		conf.Shutdown,
		middleware.Logger(conf.Log),
		middleware.Errors(conf.Log),
		middleware.Metrics(),
		middleware.Panics(), // Always has to be add the end of chain viz. Needs to be the first point of entry
	)
	// Load the routes for different versions
	v1(app, conf)
	return app
}

// v1 binds all the v1 routes
func v1(app *web.App, conf APIMuxConfig) {
	const version = "v1"
	tgh := testgrp.Handler{
		Log: conf.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
}
