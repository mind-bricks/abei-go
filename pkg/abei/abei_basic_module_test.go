package abei

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicModule(t *testing.T) {
	assert := assert.New(t)

	mID := uuid.NewString()
	m, err := NewBasicModule(
		mID,
		"root",
		"this is root module for testing",
		[]IModule{},
	)

	assert.Equal(err, nil, "error found")
	assert.NotEqual(m, nil, "nil module found")
	assert.Equal(m.GetID(), mID)
	assert.Equal(m.GetName(), "root")
}
