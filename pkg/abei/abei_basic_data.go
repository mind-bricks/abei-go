package abei

import (
	"github.com/google/uuid"
)

type BasicDataFuncCreate func(ref interface{}) (interface{}, error)
type BasicDataFuncClone func(ref interface{}) (interface{}, error)

type basicData struct {
	id    string
	cls   *basicDataClass
	value interface{}
}

func (d *basicData) GetID() string {
	return d.id
}

func (d *basicData) GetClass() IDataClass {
	return d.cls
}

func (d *basicData) GetValue() interface{} {
	return d.value
}

func (d *basicData) Clone() (IData, error) {
	value, err := d.cls.funcClone(d.value)
	if err != nil {
		return nil, err
	}

	dCloned := basicData{
		id:    uuid.New().String(),
		cls:   d.cls,
		value: value,
	}
	return &dCloned, nil
}

type basicDataClass struct {
	module     IModule
	id         string
	name       string
	document   string
	funcCreate BasicDataFuncCreate
	funcClone  BasicDataFuncClone
}

func (dc *basicDataClass) GetModule() IModule {
	return dc.module
}

func (dc *basicDataClass) GetID() string {
	return dc.id
}

func (dc *basicDataClass) GetName() string {
	return dc.name
}

func (dc *basicDataClass) GetDocument() string {
	return dc.document
}

func (dc *basicDataClass) Create(val interface{}) (IData, error) {
	value, err := dc.funcCreate(val)
	if err != nil {
		return nil, err
	}

	d := basicData{
		id:    uuid.New().String(),
		cls:   dc,
		value: value,
	}
	return &d, nil
}

func NewBasicDataClass(
	module IModule,
	id string,
	name string,
	document string,
	funcCreate BasicDataFuncCreate,
	funcClone BasicDataFuncClone,
) (IDataClass, error) {
	dc := basicDataClass{
		module:     module,
		id:         id,
		name:       name,
		document:   document,
		funcCreate: funcCreate,
		funcClone:  funcClone,
	}
	err := module.RegisterDataClass(&dc)
	return &dc, err
}
