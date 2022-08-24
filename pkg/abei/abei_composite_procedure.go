package abei

import (
	"errors"

	"github.com/google/uuid"
)

type compositeProcedureLink struct {
	node string // input/ouput node
	slot int    // input/output index
}

type compositeProcedureCache map[string]IDataParams

type compositeProcedureNodeOutput struct {
	node *compositeProcedureNode
	slot int
}

func (o *compositeProcedureNodeOutput) Run(args IDataParams, cache compositeProcedureCache) (IData, error) {
	if o.node == nil {
		data, ok := args[o.slot]
		if !ok {
			return nil, errors.New("not found in params")

		}
		return data, nil
	}

	outputs, err := o.node.Run(args, cache)
	if err != nil {
		return nil, err
	}

	output, ok := outputs[o.slot]
	if !ok {
		return nil, errors.New("not found in node outputs")
	}

	return output, nil
}

type compositeProcedureNode struct {
	id     string
	ref    IProcedure
	inputs map[int]compositeProcedureNodeOutput
}

func (n *compositeProcedureNode) Run(args IDataParams, cache compositeProcedureCache) (IDataParams, error) {
	outputs, ok := cache[n.id]
	if ok {
		return outputs, nil
	}

	// TODO: run procedure recursively
	return nil, nil
}

type compositeProcedure struct {
	id      string
	cls     *compositeProcedureClass
	outputs map[int]compositeProcedureNodeOutput
}

func (p *compositeProcedure) GetID() string {
	return p.id
}

func (p *compositeProcedure) GetClass() IProcedureClass {
	return p.cls
}

func (p *compositeProcedure) Run(args IDataParams) (IDataParams, error) {
	cache := compositeProcedureCache{}
	outputs := IDataParams{}
	for slot, o := range p.outputs {
		output, err := o.Run(args, cache)
		if err != nil {
			return nil, err
		}
		outputs[slot] = output
	}

	return outputs, nil
}

type compositeProcedureClass struct {
	module   IModule
	id       string
	name     string
	document string
	inputs   IDataClassParams
	outputs  IDataClassParams
}

func (pc *compositeProcedureClass) GetModule() IModule {
	return pc.module
}

func (pc *compositeProcedureClass) GetID() string {
	return pc.id
}

func (pc *compositeProcedureClass) GetName() string {
	return pc.name
}

func (pc *compositeProcedureClass) GetDocument() string {
	return pc.document
}

func (pc *compositeProcedureClass) GetInputs() IDataClassParams {
	return pc.inputs
}

func (pc *compositeProcedureClass) GetOutputs() IDataClassParams {
	return pc.outputs
}

func (pc *compositeProcedureClass) Create() (IProcedure, error) {
	// TODO: create sub procedure by nodes and links hold by procedure class
	// ...

	p := compositeProcedure{
		id:  uuid.New().String(),
		cls: pc,
	}
	return &p, nil
}

func NewCompositeProcedureClass(
	module IModule,
	id string,
	name string,
	document string,
	inputs map[int]IDataClass,
	outputs map[int]IDataClass,
	nodes map[string]IProcedureClass,
	links map[compositeProcedureLink]compositeProcedureLink,
) (IProcedureClass, error) {
	// TODO: check if inputs and ouputs is all resgistered in module or its dependencies
	// ...

	pc := compositeProcedureClass{
		module:   module,
		id:       id,
		name:     name,
		document: document,
		inputs:   inputs,
		outputs:  outputs,
	}

	// TODO: load nodes and links by parsing scripts
	// ...

	err := module.RegisterProcedureClass(&pc)
	return &pc, err
}
