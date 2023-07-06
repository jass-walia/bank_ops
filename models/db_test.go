package models

import (
	"testing"

	"github.com/jass-walia/bank_ops/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initalize config.
	config.Initialize("../.env")
}

// TestOpenDB tests the database connection.
func TestOpenDB(t *testing.T) {
	// Open db connection using test database.
	config.C.DBName = config.C.TestDBName
	assert.Nil(t, OpenDB())
}

// TestMigrateDB tests the database migrations.
func TestMigrateDB(t *testing.T) {
	// Open db connection using test database.
	config.C.DBName = config.C.TestDBName
	if err := OpenDB(); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, MigrateDB())
	assert.True(t, DB.Migrator().HasTable(&Account{}))
	assert.True(t, DB.Migrator().HasTable(&Transaction{}))
}
