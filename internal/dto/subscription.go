package dto

import (
	"errors"
	"time"

	"github.com/BabichevDima/subManager/internal/models"
	"github.com/google/uuid"
)

type RequestSubscription struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,gt=0"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

type ResponseSubscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrSubscriptionExists   = errors.New("subscription already exists")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidID            = errors.New("invalid ID format")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrRecordNotFound       = errors.New("subscription not found or already deleted")
)

func FromModel(sub *models.Subscription) ResponseSubscription {
	response := ResponseSubscription{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("01-2006"),
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}

	if sub.EndDate != nil {
		endDateStr := sub.EndDate.Format("01-2006")
		response.EndDate = &endDateStr
	}

	return response
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	EndDate     string `json:"end_date"`
}

type TotalCostRequest struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type TotalCostResponse struct {
	TotalCost          float64 `json:"total_cost"`
	SubscriptionsCount int     `json:"subscriptions_count"`
}
