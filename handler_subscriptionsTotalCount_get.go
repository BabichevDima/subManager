package main

import (
	"log"
	"net/http"
	"time"

	"github.com/BabichevDima/subManager/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSubscriptionsTotalCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")
	userIdStr := q.Get("user_id")
	serviceName := q.Get("service_name")

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userId format", err)
		return
	}

	if startDate == "" || endDate == "" {
		respondWithError(w, http.StatusBadRequest, "Both start_date and end_date are required", nil)
		return
	}

	if _, err := time.Parse("01-2006", startDate); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid start_date format (use MM-YYYY)", err)
		return
	}

	log.Printf("Calculating total cost: start=%s end=%s user=%s service=%s", startDate, endDate, userIdStr, serviceName)

	dbCalculateTotalCostRow, err := cfg.DB.CalculateTotalCost(r.Context(), database.CalculateTotalCostParams{
		StartDate:   startDate,
		EndDate:     endDate,
		UserId:      userId,
		ServiceName: serviceName,
	})
	if err != nil {
		log.Printf("CalculateTotalCost error: %v", err)
		respondWithError(w, http.StatusNotFound, "Failed to get Total Cost", err)
		return
	}

	response := database.CalculateTotalCostRow{
		TotalCost:          dbCalculateTotalCostRow.TotalCost,
		SubscriptionsCount: dbCalculateTotalCostRow.SubscriptionsCount,
	}
	log.Printf("Calculation result: %+v", response)
	respondWithJSON(w, http.StatusOK, response)
}
