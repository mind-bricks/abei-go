package abei

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicDataClass(t *testing.T) {
	assert := assert.New(t)

	m, _ := NewBasicModule(
		uuid.NewString(),
		"root",
		"this is root module for testing",
		[]IModule{},
	)

	dcID := uuid.NewString()
	dc, err := NewBasicDataClass(
		m,
		dcID,
		"int32",
		"int32 class for testing",
		func(ref interface{}) (interface{}, error) {
			switch ref.(type) {
			case int32:
				return ref.(int32), nil
			case uint32:
				return int32(ref.(uint32)), nil
			case int64:
				return int32(ref.(int64)), nil
			case uint64:
				return int32(ref.(uint64)), nil
			case int:
				return int32(ref.(int)), nil
			case uint:
				return int32(ref.(uint)), nil
			case nil:
				return int32(0), nil
			}
			return nil, errors.New("invalid data type")
		},
		func(ref interface{}) (interface{}, error) { return ref.(int32), nil },
	)

	assert.Equal(err, nil, "error found")
	assert.NotEqual(dc, nil, "nil data class")
	assert.Equal(dc.GetID(), dcID)
	assert.Equal(dc.GetName(), "int32")

	d, err := dc.Create("hehe")
	assert.NotEqual(err, nil, "should raise error")
	assert.Equal(d, nil, "should be nil")

	d1, err := dc.Create(nil)
	assert.Equal(err, nil, "error found")
	assert.NotEqual(d1, nil, "nil data")
	assert.Equal(d1.GetValue().(int32), int32(0), "incorrect value")

	d2, err := dc.Create(2)
	assert.Equal(err, nil, "error found")
	assert.NotEqual(d2, nil, "nil data")
	assert.Equal(d2.GetValue().(int32), int32(2), "incorrect value")

	d3, err := d2.Clone()
	assert.Equal(err, nil, "error found")
	assert.NotEqual(d3, nil, "nil data")
	assert.Equal(d3.GetValue().(int32), int32(2), "incorrect value")
}
