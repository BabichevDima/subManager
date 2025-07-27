package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSubscriptionsDelete(w http.ResponseWriter, r *http.Request) {
	subscriptionIdStr := r.PathValue("subscriptionId")

	if subscriptionIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "subscriptionId is required", nil)
		return
	}

	subscriptionId, err := uuid.Parse(subscriptionIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscriptionId format", err)
		return
	}

	log.Printf("Deleting subscription: %s", subscriptionId)
	defer func() {
		if err == nil {
			log.Printf("Deleted subscription: %s", subscriptionId)
			log.Printf("------------")
		}
	}()

	result, err := cfg.DB.DeleteSubscription(r.Context(), subscriptionId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete subscription", err)
		return
	}

	if result == 0 {
		respondWithError(w, http.StatusNotFound, "Subscription not found or already deleted", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
