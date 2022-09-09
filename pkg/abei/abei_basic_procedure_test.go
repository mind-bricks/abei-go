package abei

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicProcedure(t *testing.T) {
	assert := assert.New(t)

	m, _ := NewBasicModule(
		uuid.NewString(),
		"root",
		"this is root module for testing",
		[]IModule{},
	)

	dcID := uuid.NewString()
	dc, _ := NewBasicDataClass(
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

	pcID := uuid.NewString()
	pc, err := NewBasicProcedureClass(
		m,
		pcID,
		"integer plus",
		"plus operation for integer",
		map[int]IDataClass{0: dc, 1: dc},
		map[int]IDataClass{0: dc},
		func(inputs map[int]interface{}) (map[int]interface{}, error) {
			a, ok := inputs[0]
			if !ok {
				return nil, errors.New("missing param 0")
			}

			b, ok := inputs[1]
			if !ok {
				return nil, errors.New("missing param 1")
			}

			c := a.(int32) + b.(int32)
			return map[int]interface{}{0: c}, nil
		},
	)

	assert.Equal(err, nil, "error found")
	assert.NotEqual(pc, nil, "nil procedure class found")

	p, err := pc.Create()
	assert.Equal(err, nil, "error found")
	assert.NotEqual(p, nil, "nil procedure found")

	a, _ := dc.Create(1)
	b, _ := dc.Create(2)
	assert.NotEqual(a, nil, "nil data found")
	assert.NotEqual(b, nil, "nil data found")

	outputs, err := p.Run(map[int]IData{0: a, 1: b})
	assert.Equal(err, nil, "error found")
	assert.Equal(len(outputs), 1)
	c, _ := outputs[0]
	assert.Equal(c.GetValue().(int32), int32(3), "incorrect output")
}
