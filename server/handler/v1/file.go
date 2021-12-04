package v1

import "github.com/mazrean/Quantainer/service"

type File struct {
	checker     *Checker
	fileService service.File
}

func NewFile(checker *Checker, fileService service.File) *File {
	return &File{
		checker:     checker,
		fileService: fileService,
	}
}
