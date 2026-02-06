package activity

type Activity struct {
	Id          int
	Name        string
	Subtitle    *string
	Description string
	Type        string
	MetaLabel   *string
	DateStart   string
	DateEnd     *string
}

type EditField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ActivityRepository interface {
	GetAllOfType(activityType string) ([]Activity, error)
	Add(item Activity) (Activity, error)
	Remove(id int) error
	Edit(id int, ops []EditField) (Activity, error)
}
