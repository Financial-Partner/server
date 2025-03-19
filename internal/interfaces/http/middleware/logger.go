package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type LoggerMiddleware struct {
	log logger.Logger
}

func NewLoggerMiddleware(log logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		log: log,
	}
}

func (m *LoggerMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		var routeName string
		if route := mux.CurrentRoute(r); route != nil {
			routeName, _ = route.GetPathTemplate()
		}
		if routeName == "" {
			routeName = r.URL.Path
		}

		wrw := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		startTime := time.Now()
		m.log.WithFields(map[string]interface{}{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"route":       routeName,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
			"event":       "request_start",
		}).Infof("API Request Started: %s %s", r.Method, routeName)

		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, contextutil.RequestIDKey, requestID))

		next.ServeHTTP(wrw, r)

		duration := time.Since(startTime)
		m.log.WithFields(map[string]interface{}{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"route":       routeName,
			"status":      wrw.statusCode,
			"duration_ms": duration.Milliseconds(),
			"event":       "request_end",
		}).Infof("API Request Completed: %s %s - Status: %d - Duration: %dms",
			r.Method,
			routeName,
			wrw.statusCode,
			duration.Milliseconds(),
		)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rww *responseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}
