package main

import (
	_ "github.com/lib/pq"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BabichevDima/subManager/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func connectToBD() (*database.Queries, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	dbQueries := database.New(db)

	return dbQueries, nil
}

func main() {
	dbQueries, err := connectToBD()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Server started on localhost:8080")
	mux := http.NewServeMux()
	apiCfg := &apiConfig{
		DB: dbQueries,
	}

	mux.Handle("/app/", middlewareLog(http.StripPrefix("/app/", http.FileServer(http.Dir("./")))))

	mux.Handle("POST /api/subscriptions", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsCreate)))
	mux.Handle("GET /api/subscriptions", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsList)))
	mux.Handle("GET /api/subscriptions/{subscriptionId}", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsGet)))
	mux.Handle("PUT /api/subscriptions/{subscriptionId}", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsUpdate)))
	mux.Handle("DELETE /api/subscriptions/{subscriptionId}", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsDelete)))

	mux.Handle("GET /api/subscriptions/total", middlewareLog(http.HandlerFunc(apiCfg.handlerSubscriptionsTotalCost)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Warning: Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
