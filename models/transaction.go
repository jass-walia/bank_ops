package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction represents fields of transactions model with gorm and json tags.
type Transaction struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UUID      uuid.UUID      `gorm:"type:uuid;not null" json:"uid"`
	Narration string         `json:"narration"`
	Credit    float64        `gorm:"type:decimal(12,2)" json:"credit"`
	Debit     float64        `gorm:"type:decimal(12,2)" json:"debit"`
	AccountID uint           `gorm:"not null" json:"-"`
	Account   Account        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"account"`
}

// BeforeCreate is a gorm hook used to set uuid before insert.
func (trans *Transaction) BeforeCreate(*gorm.DB) error {
	// Set UUID
	if trans.UUID == uuid.Nil {
		trans.UUID = uuid.New()
	}
	return nil
}

// Create inserts a transaction record into database, or returns an error.
func (trans *Transaction) Create() error {
	return DB.Create(trans).Error
}

// CreateMultiple is to create multiple transaction records, or returns an error.
func (trans *Transaction) CreateMultiple(transactions []Transaction) error {
	return DB.Create(transactions).Error
}

// GetBalance returns the balance of specific account based on transactions till now.
func (trans *Transaction) GetBalance() (*float64, error) {
	var balance *float64
	result := DB.Model(&trans).Select("SUM(credit)-SUM(debit) as balance").Where("account_id = ?", trans.Account.ID).Scan(&balance)
	if result.RowsAffected == 0 && result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return balance, result.Error
}
