package algorithmia

type DataObjectType int

const (
	File DataObjectType = iota
	Directory
)

const DataObjectNone DataObjectType = -1

func (obj DataObjectType) IsFile() bool {
	return obj == File
}

func (obj DataObjectType) IsDir() bool {
	return obj == Directory
}

func (obj DataObjectType) Type() DataObjectType {
	return obj
}

type DataObject interface {
	IsFile() bool
	IsDir() bool
	Type() DataObjectType
}
