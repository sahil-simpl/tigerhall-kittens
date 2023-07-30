package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"tigerhallKittens/app"
	"tigerhallKittens/app/lib/db"
	"tigerhallKittens/app/lib/logger"
)

var ctx context.Context

const (
	cpuCount = 20
)

func init() {
	if err := app.LoadEnv(); err != nil {
		panic(err)
	}
	ctx = context.Background()
	app.SetupLogger(app.Env.Environment)
	app.SetupDBConnection(ctx)
	runtime.GOMAXPROCS(cpuCount)
	logger.I(context.Background(), "Initialization complete")
}

func main() {
	defer logger.Sync()
	defer db.Close()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:     strings.Split(app.Env.AllowedOrigins, ","),
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"*"},
	})

	router := httprouter.New()

	handler := corsMiddleware.Handler(router)
	server := &http.Server{Addr: fmt.Sprintf(":%s", app.Env.Port), Handler: handler}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.W(ctx, "Stopping server gracefully")
		}
	}()
	logger.W(ctx, fmt.Sprintf("Starting web server at port %v", app.Env.Port))

	<-done
	logger.W(ctx, fmt.Sprintf("Stopping web server at port %v", app.Env.Port))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		defer logger.Sync()
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.E(ctx, err, fmt.Sprintf("error while shutting down server:%+v", err))
	}
}
