package middleware

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"go.uber.org/zap"
)

// Authorization ensures that the request has a valid API key.
func Authorization(apiKey string, logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			serror.JSON(w, logger, serror.ErrorResponse{
				Message: "Invalid or missing API key. Please provide a valid API key.",
				Code:    http.StatusUnauthorized,
			})

			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserAgent ensures that the request has a valid user agent.
func UserAgent(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() == "" {
			serror.JSON(w, logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "User agent is missing. Please provide a valid user agent.",
			})

			return
		}

		next.ServeHTTP(w, r)
	})
}

// PrivacyPolicy adds a privacy policy header to the response.
func PrivacyPolicy(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Privacy-Policy", uri)

		next.ServeHTTP(w, r)
	})
}

// TermsOfService adds a terms of service header to the response.
func TermsOfService(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Terms-Of-Service", uri)

		next.ServeHTTP(w, r)
	})
}
