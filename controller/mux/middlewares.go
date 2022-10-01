package mux

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/config"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type middleware struct {
	cfg *config.Secrets
}

func (*handler) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\t%s\t%s\n", time.Now().Format(time.RFC3339), r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

// access to normal users
func (h *handler) authorizationMiddlewareMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken, err := getAuthToken(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		claims, err := h.srv.ClaimsFromToken(reqToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if claims != nil {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (h *handler) authorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken, err := getAuthToken(r)
		if err != nil {
			writeJson(w, http.StatusUnauthorized, dto.Error{Status: dto.StatusError, Error: err.Error()})
			return
		}
		claims, err := h.srv.ClaimsFromToken(reqToken)
		if err != nil {
			writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
			return
		}

		if claims != nil {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next(w, r.WithContext(ctx))
		}
	}
}

func getAuthToken(r *http.Request) (string, error) {
	// reqToken := w.Header.Get("Authorization")
    reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		// Error: Bearer token not in proper format
		return "", fmt.Errorf("auth header not in proper format")
	}

	reqToken = strings.TrimSpace(splitToken[1])
	return reqToken, nil
}
