package algorithmia

type DataObjectType int

const (
	File DataObjectType = iota
	Directory
)

const DataObjectNone DataObjectType = -1

//Returns whether object is a file
func (obj DataObjectType) IsFile() bool {
	return obj == File
}

//Returns whether object is a directory
func (obj DataObjectType) IsDir() bool {
	return obj == Directory
}

//Returns type of DataObject
func (obj DataObjectType) Type() DataObjectType {
	return obj
}

type DataObject interface {
	IsFile() bool         //Returns whether object is a file
	IsDir() bool          //Returns whether object is a directory
	Type() DataObjectType //Returns type of this DataObject
}
