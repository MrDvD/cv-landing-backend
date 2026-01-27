package files

import "os"

type FileHandler struct {
	BasePath []string
}

func (h *FileHandler) Get(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	return File{
		Content: string(data),
	}, nil
}
