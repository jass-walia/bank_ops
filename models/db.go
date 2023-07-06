package models

import (
	"fmt"

	"github.com/jass-walia/bank_ops/config"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB holds the database connection.
var DB *gorm.DB

// OpenDB opens a connection with the database or returns an err.
func OpenDB() (err error) {
	c := config.C

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)

	glog.V(2).Infof("Connecting to database, details: host=%s port=%d user=%s dbname=%s sslmodel=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.DBSSLMode)

	// make a connection to the database specified in the psqlInfo
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		glog.V(2).Infoln("Database is connected")
	}

	return
}

// MigrateDB migrates all the tables.
func MigrateDB() error {
	glog.V(2).Info("init database")

	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				// Account
				type Account struct {
					gorm.Model
					UUID              uuid.UUID `gorm:"type:uuid;not null"`
					AccountHolderName string    `gorm:"not null"`
					AccountNumber     uint      `gorm:"type:serial;unique"`
					AccountType       string    `gorm:"not null"`
				}
				return tx.AutoMigrate(&Account{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("accounts")
			},
		},
		{
			ID: "2",
			Migrate: func(tx *gorm.DB) error {
				type Transaction struct {
					gorm.Model
					UUID      uuid.UUID `gorm:"type:uuid;not null"`
					Narration string
					Credit    float64 `gorm:"type:decimal(12,2)"`
					Debit     float64 `gorm:"type:decimal(12,2)"`
					AccountID uint    `gorm:"not null"`
					Account   Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
				}
				return tx.AutoMigrate(&Transaction{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("transactions")
			},
		},
	})

	return m.Migrate()
}
