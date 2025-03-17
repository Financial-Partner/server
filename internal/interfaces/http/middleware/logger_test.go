package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerMiddleware(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	middleware := middleware.NewLoggerMiddleware(nopLogger)
	assert.NotNil(t, middleware)
}

func TestLoggerMiddleware_LogRequest(t *testing.T) {
	t.Run("logs request start and end with correct fields", func(t *testing.T) {
		nopLogger := logger.NewNopLogger()
		loggerMiddleware := middleware.NewLoggerMiddleware(nopLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value(contextutil.RequestIDKey)
			assert.NotNil(t, requestID)
			assert.IsType(t, "", requestID)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("response body"))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("User-Agent", "test-agent")

		router := mux.NewRouter()
		router.Handle("/test", loggerMiddleware.LogRequest(nextHandler)).Methods("GET")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("logs request with non-200 status code", func(t *testing.T) {
		nopLogger := logger.NewNopLogger()
		loggerMiddleware := middleware.NewLoggerMiddleware(nopLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request"))
		})

		req := httptest.NewRequest("GET", "/test-error", nil)

		router := mux.NewRouter()
		router.Handle("/test-error", loggerMiddleware.LogRequest(nextHandler)).Methods("GET")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("works without mux router", func(t *testing.T) {
		nopLogger := logger.NewNopLogger()
		loggerMiddleware := middleware.NewLoggerMiddleware(nopLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/no-mux", nil)
		rr := httptest.NewRecorder()

		handler := loggerMiddleware.LogRequest(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestResponseWriterWrapper(t *testing.T) {
	t.Run("captures status code correctly", func(t *testing.T) {
		nopLogger := logger.NewNopLogger()
		loggerMiddleware := middleware.NewLoggerMiddleware(nopLogger)

		statusCodes := []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError}

		for _, code := range statusCodes {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			})

			wrappedHandler := loggerMiddleware.LogRequest(handler)

			req := httptest.NewRequest("GET", "/status-test", nil)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			assert.Equal(t, code, rr.Code)
		}
	})
}
