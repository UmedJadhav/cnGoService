package middleware

import (
	"cnGoService/business/system/validate"
	"cnGoService/foundation/web"
	"context"
	"net/http"

	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				var er validate.ErrorResponse
				var status int
				switch act := validate.Cause(err).(type) {
				case validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest
				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status
				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// IF we recieve the shutdown err , we need to return it to the base handler to shutdown the service
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// The error has been handled , so we stop propogating the error further
			return nil
		}
		return h
	}
	return m
}
