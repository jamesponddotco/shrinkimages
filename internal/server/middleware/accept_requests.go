package middleware

import (
	"fmt"
	"net/http"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"go.uber.org/zap"
)

// AcceptRequests is a middleware that only allows requests with the specified
// request methods.
func AcceptRequests(methods []string, logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)

				return
			}
		}

		w.Header().Set("Allow", fmt.Sprintf("%v", methods))

		serror.JSON(w, logger, serror.ErrorResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: fmt.Sprintf("Method %s not allowed.", r.Method),
		})
	})
}
