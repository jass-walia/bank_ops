package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeConfig(t *testing.T) {
	Initialize("../.env")

	assert.NotEqual(t, C, Configurations{})
}
