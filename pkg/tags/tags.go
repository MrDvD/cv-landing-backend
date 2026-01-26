package tags

type Tag struct {
	Name       string
	Type       string
	ActivityId int
}

type TagsRepository interface {
	GetAll(activityId int) ([]Tag, error)
	GetAllOfType(activityId int, tagType string) ([]Tag, error)
}
