package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/BabichevDima/subManager/internal/database"
	"github.com/google/uuid"
)

type Subscription struct {
	ID          string       `json:"id"`
	ServiceName string       `json:"service_name"`
	Price       int32        `json:"price"`
	UserID      string       `json:"user_id"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     sql.NullTime `json:"end_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (cfg *apiConfig) handlerSubscriptionsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if params.ServiceName == "" {
		respondWithError(w, http.StatusBadRequest, "ServiceName is required", err)
		return
	}
	if params.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price must be positive", err)
		return
	}
	if params.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "UserID is required", err)
		return
	}
	if params.StartDate == "" {
		respondWithError(w, http.StatusBadRequest, "StartDate is required", err)
		return
	} else if _, err := time.Parse("01-2006", params.StartDate); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid StartDate format (use MM-YYYY)", err)
		return
	}
	if params.EndDate != "" {
		if _, err := time.Parse("01-2006", params.EndDate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid EndDate format (use MM-YYYY)", err)
			return
		}
	}

	userUUID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UserID format", err)
		return
	}

	if _, err := cfg.DB.GetSubscriptionByServiceAndUser(r.Context(),
		database.GetSubscriptionByServiceAndUserParams{
			ServiceName: params.ServiceName,
			UserID:      userUUID,
		}); err == nil {
		respondWithError(w, http.StatusConflict, "Subscription already exists", nil)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	subscription, err := cfg.DB.CreateSubscription(r.Context(), database.CreateSubscriptionParams{
		ServiceName: params.ServiceName,
		Price:       int32(params.Price),
		UserID:      userUUID,
		ToDate:      params.StartDate,
		Column5:     params.EndDate,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusConflict, "Subscription already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Subscription{
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
