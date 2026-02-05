package files

import "os"

type FileHandler struct {
	Path []string
}

func (h FileHandler) Get(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	return File{
		Content: string(data),
	}, nil
}

func (h FileHandler) BasePath() []string {
	return h.Path
}
