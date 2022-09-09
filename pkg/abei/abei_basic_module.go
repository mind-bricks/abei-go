package abei

import (
	"errors"
	"fmt"
)

type basicModule struct {
	id               string
	name             string
	document         string
	freezed          bool
	dependencies     map[string]IModule
	dataClasses      map[string]IDataClass
	procedureClasses map[string]IProcedureClass
}

func (m *basicModule) GetID() string {
	return m.id
}

func (m *basicModule) GetName() string {
	return m.name
}

func (m *basicModule) GetDocument() string {
	return m.document
}

func (m *basicModule) GetDependency(moduleID string) IModule {
	dep, ok := m.dependencies[moduleID]
	if ok {
		return dep
	}

	for _, dep := range m.dependencies {
		depSup := dep.GetDependency(moduleID)
		if depSup != nil {
			return depSup
		}
	}

	return nil
}

func (m *basicModule) GetDataClass(dataClassID string) IDataClass {
	dc, ok := m.dataClasses[dataClassID]
	if ok {
		return dc
	}
	return nil
}

func (m *basicModule) GetProcedureClass(procedureClassID string) IProcedureClass {
	pc, ok := m.procedureClasses[procedureClassID]
	if ok {
		return pc
	}

	return nil
}

func (m *basicModule) listDependencies() []IModule {
	deps := []IModule{}
	for _, v := range m.dependencies {
		deps = append(deps, v)
	}
	return deps
}

func (m *basicModule) listDataClasses() []IDataClass {
	dataClasses := []IDataClass{}
	for _, v := range m.dataClasses {
		dataClasses = append(dataClasses, v)
	}
	return dataClasses
}

func (m *basicModule) listProcedureClasses() []IProcedureClass {
	procClasses := []IProcedureClass{}
	for _, v := range m.procedureClasses {
		procClasses = append(procClasses, v)
	}
	return procClasses
}

func (m *basicModule) RegisterDataClass(dc IDataClass) error {
	if m.freezed {
		return errors.New(fmt.Sprintf("module %s freezed", m.name))
	}
	if dc.GetModule() != m {
		return errors.New(fmt.Sprintf(
			"data class %s belongs to a different module", dc.GetName()))
	}

	m.dataClasses[dc.GetID()] = dc

	return nil
}

func (m *basicModule) RegisterProcedureClass(pc IProcedureClass) error {
	if m.freezed {
		return errors.New(fmt.Sprintf("module %s freezed", m.name))
	}
	if pc.GetModule() != m {
		return errors.New(fmt.Sprintf(
			"procedure class %s belongs to a different module", pc.GetName()))
	}

	m.procedureClasses[pc.GetID()] = pc

	return nil
}

func (m *basicModule) RegisterFreeze() error {
	if len(m.procedureClasses) == 0 {
		return errors.New(fmt.Sprintf("module %s has now procedures", m.name))
	}

	m.freezed = true
	return nil
}

func NewBasicModule(
	id string,
	name string,
	document string,
	dependencies []IModule,
) (IModule, error) {
	m := basicModule{
		id:               id,
		name:             name,
		document:         document,
		freezed:          false,
		dependencies:     map[string]IModule{},
		dataClasses:      map[string]IDataClass{},
		procedureClasses: map[string]IProcedureClass{},
	}
	// TODO: check circle dependencies
	for _, dep := range dependencies {
		m.dependencies[dep.GetID()] = dep
	}
	return &m, nil
}
