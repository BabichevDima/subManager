package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSubscriptionsGet(w http.ResponseWriter, r *http.Request) {
	subscriptionIdStr := r.PathValue("subscriptionId")

	subscriptionId, err := uuid.Parse(subscriptionIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscriptionId format", err)
		return
	}

	log.Printf("Fetching subscription: %s", subscriptionIdStr)
	defer func() {
		if err == nil {
			log.Printf("Fetched subscription: %s", subscriptionId)
		} else {
			log.Printf("Error fetching subscription: %v", err)
		}
		log.Printf("------------")
	}()

	dbSubscription, err := cfg.DB.GetSubscription(r.Context(), subscriptionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Subscription not found", nil)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to get subscription", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, Subscription{
		ID:          dbSubscription.ID.String(),
		ServiceName: dbSubscription.ServiceName,
		Price:       dbSubscription.Price,
		UserID:      dbSubscription.UserID.String(),
		StartDate:   dbSubscription.StartDate,
		CreatedAt:   dbSubscription.CreatedAt,
		UpdatedAt:   dbSubscription.UpdatedAt,
		EndDate:     dbSubscription.EndDate,
	})
}
