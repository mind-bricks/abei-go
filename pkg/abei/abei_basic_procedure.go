package abei

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type BasicProcedureFuncRun func(inputs map[int]interface{}) (map[int]interface{}, error)

type basicProcedure struct {
	id  string // runing instance
	cls *basicProcedureClass
}

func (p *basicProcedure) GetID() string {
	return p.id
}

func (p *basicProcedure) GetClass() IProcedureClass {
	return p.cls
}

func (p *basicProcedure) Run(inputs IDataParams) (IDataParams, error) {
	// verify and unwrap input data
	inputClasses := p.cls.GetInputs()
	if len(inputs) != len(inputClasses) {
		return nil, errors.New(fmt.Sprintf(
			"procedure need %d inputs, but %d were given",
			len(inputClasses),
			len(inputs),
		))
	}

	inputDataValues := make(map[int]interface{})
	for slot, inputClass := range inputClasses {
		inputData, ok := inputs[slot]
		if !ok {
			return nil, errors.New(fmt.Sprintf(
				"procedure missing input class %s at slot %d",
				inputClass.GetID(),
				slot,
			))
		}

		if inputData.GetClass().GetID() != inputClass.GetID() {
			return nil, errors.New(fmt.Sprintf(
				"procedure need input class %s at slot %d, where class %s were given",
				inputClass.GetID(),
				slot,
				inputData.GetClass().GetID(),
			))
		}

		inputDataValues[slot] = inputData.GetValue()
	}

	// run builtin function
	outputDataValues, err := p.cls.funcRun(inputDataValues)
	if err != nil {
		return nil, err
	}

	// wrap output data
	outputClasses := p.cls.GetOutputs()
	for len(outputDataValues) != len(outputClasses) {
		return nil, errors.New(fmt.Sprintf(
			"procedure should return %d outputs instead of %d",
			len(outputClasses),
			len(outputDataValues)))
	}

	outputDatas := make(map[int]IData)
	for slot, outputClass := range outputClasses {
		outputDataValue, ok := outputDataValues[slot]
		if !ok {
			return nil, errors.New(fmt.Sprintf(
				"procedure missing output class %s at slot %d",
				outputClass.GetID(), slot))
		}

		// create data to hold value
		outputData, err := outputClass.Create(outputDataValue)
		if err != nil {
			return nil, err
		}

		outputDatas[slot] = outputData
	}

	return outputDatas, nil
}

type basicProcedureClass struct {
	module   IModule
	id       string
	name     string
	version  string
	document string
	inputs   IDataClassParams
	outputs  IDataClassParams
	funcRun  BasicProcedureFuncRun
}

func (dc *basicProcedureClass) GetModule() IModule {
	return dc.module
}

func (pc *basicProcedureClass) GetID() string {
	return pc.id
}

func (pc *basicProcedureClass) GetName() string {
	return pc.name
}

func (pc *basicProcedureClass) GetDocument() string {
	return pc.document
}

func (pc *basicProcedureClass) GetInputs() IDataClassParams {
	return pc.inputs
}

func (pc *basicProcedureClass) GetOutputs() IDataClassParams {
	return pc.outputs
}

func (pc *basicProcedureClass) Create() (IProcedure, error) {
	p := basicProcedure{
		id:  uuid.New().String(),
		cls: pc,
	}
	return &p, nil
}

func NewBasicProcedureClass(
	module IModule,
	id string,
	name string,
	document string,
	inputs map[int]IDataClass,
	outputs map[int]IDataClass,
	funcRunc BasicProcedureFuncRun,
) (IProcedureClass, error) {

	moduleID := module.GetID()

	for _, dc := range inputs {
		module := dc.GetModule()
		if module.GetID() != moduleID && module.GetDependency(moduleID) == nil {
			return nil, errors.New("input data class not found in module or its dependecies")
		}
	}

	for _, dc := range outputs {
		module := dc.GetModule()
		if module.GetID() != moduleID && module.GetDependency(moduleID) == nil {
			return nil, errors.New("ouput data class not found in module or its dependecies")
		}
	}

	pc := basicProcedureClass{
		module:   module,
		id:       id,
		name:     name,
		document: document,
		inputs:   inputs,
		outputs:  outputs,
		funcRun:  funcRunc,
	}
	err := module.RegisterProcedureClass(&pc)
	return &pc, err
}
