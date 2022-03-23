package model

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type Default struct {
	ID        uuid.UUID  	`gorm:"primaryKey;"`
	IsDeleted bool  		`gorm:"not null;default:0;"`
	CreatedAt time.Time 	`gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time 	`gorm:"autoUpdateTime:milli;not null;"`
	CreatedBy uint  		`gorm:"not null;default:0;"`
	UpdatedBy uint   		`gorm:"not null;default:0;"`
	CurrentVersion uint  	`gorm:"not null;default:0;"`
}