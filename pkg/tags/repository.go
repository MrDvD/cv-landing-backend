package tags

import "database/sql"

type TagsHandler struct {
	DB *sql.DB
}

func (h *TagsHandler) GetAll(activityId int) ([]Tag, error) {
	return []Tag{}, nil
}

func (h *TagsHandler) GetAllOfType(activityId int, tagType string) ([]Tag, error) {
	return []Tag{}, nil
}
