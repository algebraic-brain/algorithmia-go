package algorithmia

type DataObject int

const (
	File DataObject = iota
	Directory
)

func (obj DataObject) IsFile() bool {
	return obj == File
}

func (obj DataObject) IsDir() bool {
	return obj == Directory
}
