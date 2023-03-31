package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair to the application server mux
func (a *App) Handle(method string, group string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// pre code processing
		if err := handler(r.Context(), w, r); err != nil {
			return
		}
		// post code processing
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	a.ContextMux.Handle(method, finalPath, h)
}
