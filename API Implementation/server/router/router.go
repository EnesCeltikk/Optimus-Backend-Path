package router

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Router struct {
	mux         *http.ServeMux
	limiter     *rate.Limiter
	middlewares []func(http.Handler) http.Handler
}

type RouteGroup struct {
	prefix      string
	middlewares []func(http.Handler) http.Handler
	router      *Router
}

func NewRouter() *Router {
	return &Router{
		mux:     http.NewServeMux(),
		limiter: rate.NewLimiter(rate.Every(time.Second), 100),
	}
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		middlewares: make([]func(http.Handler) http.Handler, 0),
		router:      r,
	}
}

func (g *RouteGroup) Use(middleware func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, middleware)
}

func (g *RouteGroup) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	// Combine group middlewares with route-specific middlewares
	allMiddlewares := append(g.middlewares, middlewares...)
	
	// Create the final handler with all middlewares
	var finalHandler http.Handler = handler
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		finalHandler = allMiddlewares[i](finalHandler)
	}

	g.router.mux.Handle(g.prefix+pattern, finalHandler)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	// Create the final handler with all middlewares
	var finalHandler http.Handler = handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}

	r.mux.Handle(pattern, finalHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// Handle preflight requests
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Rate limiting
	if !r.limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Request logging
	start := time.Now()
	r.mux.ServeHTTP(w, req)
	log.Info().
		Str("method", req.Method).
		Str("path", req.URL.Path).
		Dur("duration", time.Since(start)).
		Msg("Request processed")
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
} 