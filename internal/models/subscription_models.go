package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ServiceName string     `gorm:"type:varchar(100);not null"`
	Price       int        `gorm:"type:integer;not null"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	StartDate   time.Time  `gorm:"type:date;not null"`
	EndDate     *time.Time `gorm:"type:date;null"`
	CreatedAt   time.Time  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;not null;default:now()"`
}
