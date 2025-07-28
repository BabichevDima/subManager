package usecase

import (
	"time"

	"github.com/BabichevDima/subManager/internal/dto"
	"github.com/BabichevDima/subManager/internal/models"
	"github.com/BabichevDima/subManager/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionUsecase struct {
	repo *repository.SubscriptionRepository
}

func NewUserUsecase(r *repository.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: r}
}

func (u *SubscriptionUsecase) Subscribe(request dto.RequestSubscription) (dto.ResponseSubscription, error) {
	startDate, err := time.Parse("01-2006", request.StartDate)
	if err != nil {
		return dto.ResponseSubscription{}, err
	}

	var endDate *time.Time
	if request.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", request.EndDate)
		if err != nil {
			return dto.ResponseSubscription{}, err
		}
		endDate = &parsedEndDate
	}

	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		return dto.ResponseSubscription{}, err
	}

	exists, err := u.repo.Exists(request.ServiceName, userUUID)
	if err != nil {
		return dto.ResponseSubscription{}, err
	}
	if exists {
		return dto.ResponseSubscription{}, dto.ErrSubscriptionExists
	}

	resp := &models.Subscription{
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := u.repo.Create(resp); err != nil {
		return dto.ResponseSubscription{}, err
	}

	return dto.ResponseSubscription{
		ID:          resp.ID,
		ServiceName: resp.ServiceName,
		Price:       resp.Price,
		UserID:      resp.UserID,
		StartDate:   resp.StartDate.Format("01-2006"),
		EndDate:     formatTimePtr(resp.EndDate),
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format("01-2006")
	return &formatted
}

func (u *SubscriptionUsecase) GetSubscriptionByID(id string) (dto.ResponseSubscription, error) {
	subscriptionId, err := uuid.Parse(id)
	if err != nil {
		return dto.ResponseSubscription{}, dto.ErrInvalidID
	}

	subscriptionResp, err := u.repo.GetByID(subscriptionId)
	if err != nil {
		return dto.ResponseSubscription{}, err
	}

	return dto.ResponseSubscription{
		ID:          subscriptionResp.ID,
		ServiceName: subscriptionResp.ServiceName,
		Price:       subscriptionResp.Price,
		UserID:      subscriptionResp.UserID,
		StartDate:   subscriptionResp.StartDate.Format("01-2006"),
		EndDate:     formatTimePtr(subscriptionResp.EndDate),
		CreatedAt:   subscriptionResp.CreatedAt,
		UpdatedAt:   subscriptionResp.UpdatedAt,
	}, nil
}

func (u *SubscriptionUsecase) GetAllSubscriptions(page, pageSize int) ([]dto.ResponseSubscription, int64, error) {
	var subscriptions []models.Subscription
	var total int64

	offset := (page - 1) * pageSize

	subscriptions, total, err := u.repo.GetAll(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ResponseSubscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = dto.FromModel(&sub)
	}

	return result, total, nil
}

func (u *SubscriptionUsecase) DeleteSubscription(id string) error {
	subscriptionId, err := uuid.Parse(id)
	if err != nil {
		return dto.ErrInvalidID
	}

	return u.repo.Delete(subscriptionId)
}

func (u *SubscriptionUsecase) UpdateSubscription(id string, req dto.UpdateSubscriptionRequest) (dto.ResponseSubscription, error) {
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return dto.ResponseSubscription{}, dto.ErrInvalidID
	}

	existing, err := u.repo.GetByID(subscriptionID)
	if err != nil {
		return dto.ResponseSubscription{}, dto.ErrRecordNotFound
	}

	if req.ServiceName != "" {
		existing.ServiceName = req.ServiceName
	}

	if req.Price > 0 {
		existing.Price = req.Price
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return dto.ResponseSubscription{}, dto.ErrInvalidFormat
		}
		endDate = &parsedEndDate
	}
	existing.EndDate = endDate

	if err := u.repo.Update(existing); err != nil {
		return dto.ResponseSubscription{}, err
	}

	return dto.ResponseSubscription{
		ID:          existing.ID,
		ServiceName: existing.ServiceName,
		Price:       existing.Price,
		UserID:      existing.UserID,
		StartDate:   existing.StartDate.Format("01-2006"),
		EndDate:     formatTimePtr(existing.EndDate),
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}, nil
}

func (u *SubscriptionUsecase) CalculateTotalCost(req dto.TotalCostRequest) (dto.TotalCostResponse, error) {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return dto.TotalCostResponse{}, dto.ErrInvalidID
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return dto.TotalCostResponse{}, err
	}

	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		return dto.TotalCostResponse{}, err
	}

	if endDate.Before(startDate) {
		return dto.TotalCostResponse{}, err
	}

	total, count, err := u.repo.CalculateTotalCost(userUUID, req.ServiceName, startDate, endDate)
	if err != nil {
		return dto.TotalCostResponse{}, err
	}

	return dto.TotalCostResponse{
		TotalCost:          total,
		SubscriptionsCount: count,
	}, nil
}
