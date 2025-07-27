package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerSubscriptionsList(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching subscriptions list")

	subscriptionList, err := cfg.DB.GetSubscriptionList(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get subscriptions list", err)
		return
	} else if len(subscriptionList) == 0 {
		respondWithJSON(w, http.StatusOK, []Subscription{})
		return
	}

	subscriptions := make([]Subscription, len(subscriptionList))

	defer func() {
		if err == nil {
			log.Printf("Returned %d subscriptions list", len(subscriptions))
			log.Printf("------------")
		}
	}()

	for i, subscription := range subscriptionList {
		subscriptions[i] = Subscription{
			ID:          subscription.ID.String(),
			ServiceName: subscription.ServiceName,
			Price:       subscription.Price,
			UserID:      subscription.UserID.String(),
			StartDate:   subscription.StartDate,
			CreatedAt:   subscription.CreatedAt,
			UpdatedAt:   subscription.UpdatedAt,
			EndDate:     subscription.EndDate,
		}
	}

	respondWithJSON(w, http.StatusOK, subscriptions)
}
