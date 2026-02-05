package tags

type Tag struct {
	Id         int
	Name       string
	Type       string
	ActivityId int
	Priority   *int
}

type TagFilter struct {
	ActivityID *int
	TagType    *string
}

type TagsRepository interface {
	Get(filter TagFilter) ([]Tag, error)
	Add(item Tag) (Tag, error)
	Remove(id int) error
}
