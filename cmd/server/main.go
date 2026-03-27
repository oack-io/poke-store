package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oack-io/poke-store/internal/data"
	"github.com/oack-io/poke-store/internal/handler"
	"github.com/oack-io/poke-store/internal/middleware"
	"github.com/oack-io/poke-store/internal/store"
)

var (
	version   = "dev"
	commitSHA = "unknown"
	buildTime = "unknown"
)

//go:embed all:static
var staticFiles embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":6001"
	}

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "poke-store-dev-secret-change-me"
	}

	catalog := data.NewCatalog()
	sessions := store.NewSessionStore(secret)
	carts := store.NewCartStore()

	authMW := middleware.Auth(sessions)

	mux := http.NewServeMux()

	// API routes
	authHandler := handler.NewAuth(sessions, carts)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/logout", authMW(authHandler.Logout))
	mux.HandleFunc("GET /api/me", authMW(authHandler.Me))

	storeHandler := handler.NewStore(catalog)
	mux.HandleFunc("GET /api/pokemon", storeHandler.List)
	mux.HandleFunc("GET /api/pokemon/{id}", storeHandler.Get)

	cartHandler := handler.NewCart(carts, catalog)
	mux.HandleFunc("GET /api/cart", authMW(cartHandler.Get))
	mux.HandleFunc("POST /api/cart/add", authMW(cartHandler.Add))
	mux.HandleFunc("POST /api/cart/remove", authMW(cartHandler.Remove))
	mux.HandleFunc("POST /api/cart/clear", authMW(cartHandler.Clear))

	checkoutHandler := handler.NewCheckout(carts, catalog)
	mux.HandleFunc("POST /api/checkout", authMW(checkoutHandler.Process))

	mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		handler.JSON(w, http.StatusOK, map[string]string{
			"version":   version,
			"commitSHA": commitSHA,
			"buildTime": buildTime,
		})
	})

	// Static files (Astro build output)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		slog.Error("mount static fs", "scope", "lifecycle", "error", err)
		os.Exit(1)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("GET /", fileServer)

	srv := &http.Server{
		Addr:         addr,
		Handler:      middleware.Logger(middleware.Recovery(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting",
			"scope", "lifecycle",
			"addr", addr,
			"version", version,
			"commit", commitSHA,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "scope", "lifecycle", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server shutting down", "scope", "lifecycle")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "scope", "lifecycle", "error", err)
	}
	slog.Info("server stopped", "scope", "lifecycle")
}
