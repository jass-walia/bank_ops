package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account represents fields of accounts model with gorm and json tags.
type Account struct {
	ID                uint           `gorm:"primarykey" json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	UUID              uuid.UUID      `gorm:"type:uuid;not null" json:"uid"`
	AccountHolderName string         `gorm:"not null" json:"account_holder_name"`
	AccountNumber     *uint          `gorm:"default:nextval('accounts_account_number_seq')" json:"account_number"`
	AccountType       string         `gorm:"not null" json:"account_type"`
}

// BeforeCreate is a gorm hook used to set uuid before insert.
func (acc *Account) BeforeCreate(*gorm.DB) error {
	// Set UUID
	if acc.UUID == uuid.Nil {
		acc.UUID = uuid.New()
	}
	return nil
}

// Create inserts an account record into database, or returns an error.
func (acc *Account) Create() error {
	return DB.Create(acc).Error
}

// GetByUUID returns an account record for given uuid or an error.
func (acc *Account) GetByUUID(UUID string) (*Account, error) {
	result := DB.Where("uuid = ?", UUID).First(&acc)
	if result.RowsAffected == 0 && result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acc, result.Error
}
