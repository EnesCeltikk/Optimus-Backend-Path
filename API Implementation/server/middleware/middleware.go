package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"CoreImplementation/services"

	"github.com/rs/zerolog/log"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := services.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("user_role").(string)
			
			hasRole := false
			for _, role := range roles {
				if role == userRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidationMiddleware(validator func(interface{}) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestBody interface{}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if err := validator(requestBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Store validated body in context
			ctx := context.WithValue(r.Context(), "validated_body", requestBody)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Interface("error", err).Msg("Panic recovered")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func PerformanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:    http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// Log performance metrics
		duration := time.Since(start)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.statusCode).
			Dur("duration", duration).
			Msg("Request completed")

		// TODO: Add metrics collection (e.g., Prometheus)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLoggerMiddleware logs detailed request information
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Request started")

		next.ServeHTTP(w, r)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Msg("Request completed")
	})
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(requestsPerSecond int) func(http.Handler) http.Handler {
	limiter := make(chan struct{}, requestsPerSecond)
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	
	go func() {
		for range ticker.C {
			select {
			case limiter <- struct{}{}:
			default:
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-limiter:
				next.ServeHTTP(w, r)
			default:
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}
		})
	}
} 