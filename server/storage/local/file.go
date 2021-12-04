package local

import "fmt"

type File struct {
	fileRootPath     string
	directoryManager *DirectoryManager
}

func NewFile(directoryManager *DirectoryManager) (*File, error) {
	fileRootPath, err := directoryManager.setupDirectory("files")
	if err != nil {
		return nil, fmt.Errorf("failed to setup directory: %w", err)
	}

	return &File{
		fileRootPath:     fileRootPath,
		directoryManager: directoryManager,
	}, nil
}
