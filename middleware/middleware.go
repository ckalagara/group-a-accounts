package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ckalagara/group-a-accounts/common"
	"github.com/google/uuid"
)

func NewMiddlewareChain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqID := r.Header.Get(common.HttpHdrKeyRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set(common.HttpHdrKeyRequestID, reqID)
		}
		ctx := context.WithValue(r.Context(), common.HttpHdrKeyRequestID, reqID)
		r = r.WithContext(ctx)

		rt := time.Now()
		crw := CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		defer func() {
			duration := time.Since(rt)
			log.Printf("%d | %s %s %s %v", crw.StatusCode, reqID, r.Method, r.URL, duration)
		}()
		log.Printf("000 | %s %s %s", reqID, r.Method, r.URL)
		next.ServeHTTP(&crw, r)

	})
}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(common.HttpHdrKeyAuthorization)
		if token == "" {
			http.Error(w, common.HttpResUnauthorized, http.StatusUnauthorized)
			return
		}

		// Todo: JWT validation

		// Add audience to context
		ctx := context.WithValue(r.Context(), common.JWTKeyAudience, "shop-regionx-sectorx-idx")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.StatusCode = code
	crw.ResponseWriter.WriteHeader(code)
}
