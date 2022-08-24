package abei

type IData interface {
	GetID() string
	GetClass() IDataClass
	GetValue() interface{}
	Clone() (IData, error)
}

type IDataClass interface {
	GetModule() IModule
	GetID() string
	GetName() string
	GetDocument() string
	Create(val interface{}) (IData, error)
}

type IDataParams map[int]IData
type IDataClassParams map[int]IDataClass

type IProcedure interface {
	GetID() string
	GetClass() IProcedureClass
	Run(args IDataParams) (IDataParams, error)
}

type IProcedureClass interface {
	GetModule() IModule
	GetID() string
	GetName() string
	// GetVersion() string
	GetDocument() string
	GetInputs() IDataClassParams
	GetOutputs() IDataClassParams
	Create() (IProcedure, error)
}

type IModule interface {
	GetID() string
	GetName() string
	GetDocument() string
	GetDependency(moduleID string) IModule
	GetDataClass(dataClassID string) IDataClass
	GetProcedureClass(procedureID string) IProcedureClass
	listDependencies() []IModule
	listDataClasses() []IDataClass
	listProcedureClasses() []IProcedureClass
	RegisterDataClass(dc IDataClass) error
	RegisterProcedureClass(pc IProcedureClass) error
	RegisterFreeze() error // check and freeze module
}
