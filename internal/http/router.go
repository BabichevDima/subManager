package router

import (
	"net/http"

	_ "github.com/BabichevDima/subManager/internal/docs"
	"github.com/BabichevDima/subManager/internal/http/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, subscriptionHandler *handlers.SubscriptionHandler) {
	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	// Static files
	mux.Handle("/", http.FileServer(http.Dir("./app")))
	// API routes
	mux.Handle("POST /api/subscriptions", http.HandlerFunc(subscriptionHandler.Subscribe))
	mux.Handle("GET /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.GetSubscriptionByID))
	mux.Handle("GET /api/subscriptions", http.HandlerFunc(subscriptionHandler.GetAllSubscriptions))
	mux.Handle("DELETE /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.DeleteSubscription))
	mux.Handle("PUT /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.UpdateSubscription))
	mux.Handle("GET /api/subscriptions/total", http.HandlerFunc(subscriptionHandler.CalculateSubscriptionsCost))
}