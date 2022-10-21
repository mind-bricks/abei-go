package abei

import (
	"errors"

	"github.com/google/uuid"
)

type compositeProcedureCache map[string]IDataParams

type compositeProcedureSource struct {
	srcNode *compositeProcedureNode
	srcSlot int
}

func (o *compositeProcedureSource) Run(args IDataParams, cache compositeProcedureCache) (IData, error) {
	if o.srcNode == nil {
		data, ok := args[o.srcSlot]
		if !ok {
			return nil, errors.New("not found in params")

		}
		return data, nil
	}

	outputs, err := o.srcNode.Run(args, cache)
	if err != nil {
		return nil, err
	}

	output, ok := outputs[o.srcSlot]
	if !ok {
		return nil, errors.New("not found in node outputs")
	}

	return output, nil
}

type compositeProcedureNode struct {
	id  string
	ref IProcedure
	// map of sources where keys are targets of links or inputs of referencing procedure
	sources map[int]compositeProcedureSource
}

func (n *compositeProcedureNode) Run(args IDataParams, cache compositeProcedureCache) (IDataParams, error) {
	outputs, ok := cache[n.id]
	if ok {
		return outputs, nil
	}

	curArgs := IDataParams{}
	for slot, source := range n.sources {
		arg, err := source.Run(args, cache)
		if err != nil {
			return nil, err
		}
		curArgs[slot] = arg
	}

	outputs, err := n.ref.Run(curArgs)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

type compositeProcedure struct {
	id      string
	cls     *compositeProcedureClass
	sources map[int]compositeProcedureSource
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
	for slot, source := range p.sources {
		output, err := source.Run(args, cache)
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
	sources  []compositeProcedureSource
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

type ProcedureJoint struct {
	node string // input/ouput node
	slot int    // input/output index
}

func NewCompositeProcedureClass(
	module IModule,
	id string,
	name string,
	document string,
	inputs map[int]IDataClass,
	outputs map[int]IDataClass,
	nodes map[string]IProcedureClass,
	// links maps where keys are targets and values are sources
	links map[ProcedureJoint]ProcedureJoint,
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

	// TODO: load nodes and links
	// ...

	err := module.RegisterProcedureClass(&pc)
	return &pc, err
}
