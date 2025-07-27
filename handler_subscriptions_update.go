package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BabichevDima/subManager/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSubscriptionsUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		EndDate     string `json:"end_date"`
	}

	subscriptionIdStr := r.PathValue("subscriptionId")
	subscriptionId, err := uuid.Parse(subscriptionIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscriptionId format", err)
		return
	}

	params := parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if params.ServiceName == "" {
		respondWithError(w, http.StatusBadRequest, "ServiceName are required", errors.New("serviceName are required"))
		return
	}
	if params.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price is required and cannot be 0", errors.New("price is required and cannot be 0"))
		return
	}

	dbParams := database.UpdateSubscriptionParams{
		ID:          subscriptionId,
		ServiceName: params.ServiceName,
		Price:       int32(params.Price),
	}

	if params.EndDate != "" {
		_, err := time.Parse("01-2006", params.EndDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid EndDate format (use MM-YYYY)", err)
			return
		}
		dbParams.EndDate = params.EndDate
	} else if params.EndDate == "" {
		dbParams.EndDate = sql.NullTime{Valid: false}
	}

	log.Printf("Updating subscription %s with params: %+v", subscriptionId, params)
	defer func() {
		if err == nil {
			log.Printf("Subscription %s updated successfully", subscriptionId)
			log.Printf("------------")
		}
	}()

	subscription, err := cfg.DB.UpdateSubscription(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update subscription", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Subscription{
		ID:          subscription.ID.String(),
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID.String(),
		StartDate:   subscription.StartDate,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
		EndDate:     subscription.EndDate,
	})
}
