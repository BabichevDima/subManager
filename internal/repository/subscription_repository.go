package repository

import (
	"time"

	"github.com/BabichevDima/subManager/internal/dto"
	"github.com/BabichevDima/subManager/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) Exists(serviceName string, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).
		Where("service_name = ? AND user_id = ?", serviceName, userID).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	return &subscription, err
}

func (r *SubscriptionRepository) GetAll(offset, limit int) ([]models.Subscription, int64, error) {
	var subscriptions []models.Subscription
	var total int64

	if err := r.db.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.ErrRecordNotFound
	}

	return nil
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	result := r.db.Model(subscription).Updates(map[string]interface{}{
		"service_name": subscription.ServiceName,
		"price":        subscription.Price,
		"end_date":     subscription.EndDate,
		"updated_at":   time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *SubscriptionRepository) CalculateTotalCost(userID uuid.UUID, serviceName string, startDate, endDate time.Time) (float64, int, error) {
    var result struct {
        TotalCost float64
        Count     int
    }

    query := r.db.Model(&models.Subscription{}).
        Select("SUM(price) as total_cost, COUNT(*) as count").
        Where("user_id = ?", userID).
        Where("start_date <= ?", endDate).
        Where("(end_date IS NULL OR end_date >= ?)", startDate)

    if serviceName != "" {
        query = query.Where("service_name = ?", serviceName)
    }

    err := query.Scan(&result).Error
    if err != nil {
        return 0, 0, err
    }

    return result.TotalCost, result.Count, nil
}