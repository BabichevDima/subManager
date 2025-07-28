package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/BabichevDima/subManager/internal/dto"
	"github.com/BabichevDima/subManager/internal/http/response"
	"github.com/BabichevDima/subManager/internal/usecase"
)

type SubscriptionHandler struct {
	usecase *usecase.SubscriptionUsecase
}

func NewUserHandler(u *usecase.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: u}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := dto.RequestSubscription{}
	err := decoder.Decode(&request)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if request.ServiceName == "" {
		response.RespondWithError(w, http.StatusBadRequest, "ServiceName is required", err)
		return
	}
	if request.Price <= 0 {
		response.RespondWithError(w, http.StatusBadRequest, "Price must be positive", err)
		return
	}
	if request.UserID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "UserID is required", err)
		return
	}
	if request.StartDate == "" {
		response.RespondWithError(w, http.StatusBadRequest, "StartDate is required", err)
		return
	}

	subscriptionResponse, err := h.usecase.Subscribe(request)
	if err != nil {
		response.RespondWithError(w, http.StatusConflict, "Subscription already exists!", err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, subscriptionResponse)
}

func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	subscriptionIdStr := r.PathValue("subscriptionId")

	responseData, err := h.usecase.GetSubscriptionByID(subscriptionIdStr)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, responseData)
}

func (h *SubscriptionHandler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	subscriptions, total, err := h.usecase.GetAllSubscriptions(page, pageSize)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to get subscriptions", err)
		return
	}

	responseData := map[string]interface{}{
		"data": subscriptions,
		"pagination": map[string]interface{}{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	response.RespondWithJSON(w, http.StatusOK, responseData)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionId := r.PathValue("subscriptionId")

	if subscriptionId == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Subscription ID is required", nil)
		return
	}

	err := h.usecase.DeleteSubscription(subscriptionId)
	if err != nil {
		switch err {
		case dto.ErrRecordNotFound:
			response.RespondWithError(w, http.StatusNotFound, "Subscription not found", err)
		case dto.ErrInvalidID:
			response.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID format", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to delete subscription", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionId := r.PathValue("subscriptionId")

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.ServiceName == "" && req.Price == 0 && req.EndDate == "" {
		response.RespondWithError(w, http.StatusBadRequest, "At least one field must be provided", nil)
		return
	}

	subscriptionResponse, err := h.usecase.UpdateSubscription(subscriptionId, req)
	if err != nil {
		switch err {
		case dto.ErrRecordNotFound:
			response.RespondWithError(w, http.StatusNotFound, "Subscription not found", err)
		case dto.ErrInvalidID:
			response.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID format", err)
		case dto.ErrInvalidFormat:
			response.RespondWithError(w, http.StatusBadRequest, "Invalid end_date format. Use MM-YYYY (e.g. '12-2025')", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to update subscription", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusOK, subscriptionResponse)
}

func (h *SubscriptionHandler) CalculateSubscriptionsCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := dto.TotalCostRequest{
		UserID:      q.Get("user_id"),
		ServiceName: q.Get("service_name"),
		StartDate:   q.Get("start_date"),
		EndDate:     q.Get("end_date"),
	}

	if req.UserID == "" || req.StartDate == "" || req.EndDate == "" {
		response.RespondWithError(w, http.StatusBadRequest, "user_id, start_date and end_date are required", nil)
		return
	}

	responseData, err := h.usecase.CalculateTotalCost(req)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrInvalidID):
			response.RespondWithError(w, http.StatusBadRequest, "Invalid user_id format", err)
		case strings.Contains(err.Error(), "invalid date format"):
			response.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to calculate total cost", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusOK, responseData)
}
