package tags

type Tag struct {
	Name       string
	Type       string
	ActivityId int
}

type TagFilter struct {
	ActivityID *int
	TagType    *string
}

type TagsRepository interface {
	Get(filter TagFilter) ([]Tag, error)
}
