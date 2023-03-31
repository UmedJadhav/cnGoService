package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
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
