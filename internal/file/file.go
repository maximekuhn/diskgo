package file

type File struct {
	Name string
	Data []byte
}

func NewFile(name string, data []byte) *File {
	return &File{
		Name: name,
		Data: data,
	}
}
