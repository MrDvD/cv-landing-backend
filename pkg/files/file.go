package files

type File struct {
	Content string
}

type FileRepository interface {
	Get(path string) (File, error)
}
