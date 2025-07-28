package router

import (
	"net/http"

	"github.com/BabichevDima/subManager/internal/http/handlers"
)

func RegisterRoutes(mux *http.ServeMux, subscriptionHandler *handlers.SubscriptionHandler) {
	mux.Handle("/", http.FileServer(http.Dir("./app"))) 

	mux.Handle("POST /api/subscriptions", http.HandlerFunc(subscriptionHandler.Subscribe))
	mux.Handle("GET /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.GetSubscriptionByID))
	mux.Handle("GET /api/subscriptions", http.HandlerFunc(subscriptionHandler.GetAllSubscriptions))
	mux.Handle("DELETE /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.DeleteSubscription))
	mux.Handle("PUT /api/subscriptions/{subscriptionId}", http.HandlerFunc(subscriptionHandler.UpdateSubscription))
	mux.Handle("GET /api/subscriptions/total", http.HandlerFunc(subscriptionHandler.CalculateSubscriptionsCost))
}