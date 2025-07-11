package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"CoreImplementation/server/handlers"
	"CoreImplementation/server/middleware"
	"CoreImplementation/server/router"
	"CoreImplementation/services"

	"github.com/rs/zerolog/log"
)

func StartServer(port string, us *services.UserService, ts *services.TransactionService, bs *services.BalanceService) {
	r := router.NewRouter()
	h := handlers.NewHandler(us, ts, bs)

	r.Use(middleware.ErrorHandlerMiddleware)
	r.Use(middleware.RequestLoggerMiddleware)
	r.Use(middleware.PerformanceMiddleware)
	r.Use(middleware.RateLimitMiddleware(100)) // 100 requests per second

	r.HandleFunc("/api/v1/auth/register", h.Register)
	r.HandleFunc("/api/v1/auth/login", h.Login)

	authGroup := r.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.HandleFunc("/users", h.GetUsers, middleware.RoleMiddleware("admin"))
	authGroup.HandleFunc("/users/", h.GetUser)
	authGroup.HandleFunc("/users/update", h.UpdateUser)
	authGroup.HandleFunc("/users/delete", h.DeleteUser, middleware.RoleMiddleware("admin"))

	authGroup.HandleFunc("/transactions/credit", h.Credit)
	authGroup.HandleFunc("/transactions/debit", h.Debit)
	authGroup.HandleFunc("/transactions/transfer", h.Transfer)
	authGroup.HandleFunc("/transactions/history", h.GetTransactionHistory)
	authGroup.HandleFunc("/transactions/", h.GetTransaction)

	authGroup.HandleFunc("/balances/current", h.GetCurrentBalance)
	authGroup.HandleFunc("/balances/historical", h.GetHistoricalBalance)
	authGroup.HandleFunc("/balances/at-time", h.GetBalanceAtTime)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("Server started at http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Graceful shutdown failed")
	}

	log.Info().Msg("Server exited")
}
