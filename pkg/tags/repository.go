package tags

import (
	"database/sql"
	"fmt"
	"strings"
)

type TagsHandler struct {
	DB *sql.DB
}

func (h *TagsHandler) Get(filter TagFilter) ([]Tag, error) {
	query, whereValues := buildGetQuery(filter)
	rows, err := h.DB.Query(query, whereValues...)
	if err != nil {
		return []Tag{}, nil
	}
	tags := []Tag{}
	for rows.Next() {
		raw := &struct {
			name       string
			tagType    string
			activityId int
			priority   *int
		}{}
		err := rows.Scan(&raw.name, &raw.tagType, &raw.activityId, &raw.priority)
		if err != nil {
			return []Tag{}, err
		}
		tags = append(tags, Tag{
			Name:       raw.name,
			Type:       raw.tagType,
			ActivityId: raw.activityId,
		})
	}
	return tags, nil
}

func buildGetQuery(filter TagFilter) (string, []any) {
	baseQuery := "select name, type, activity_id, priority from TAGS"
	whereConditions := []string{}
	whereValues := []any{}
	if filter.ActivityID != nil {
		whereConditions = append(whereConditions, "activity_id = $")
		whereValues = append(whereValues, *filter.ActivityID)
	}
	if filter.TagType != nil {
		whereConditions = append(whereConditions, "type = $")
		whereValues = append(whereValues, *filter.TagType)
	}
	var query string
	if len(whereConditions) != 0 {
		for i, val := range whereConditions {
			whereConditions[i] = fmt.Sprintf("%s%d", val, i+1)
		}
		query = fmt.Sprintf("%s where %s", baseQuery, strings.Join(whereConditions, " and "))
	} else {
		query = baseQuery
	}
	return query, whereValues
}

func (h *TagsHandler) Add(item Tag) (Tag, error) {
	var tagId int
	err := h.DB.QueryRow("insert into TAGS(name, type, activity_id, priority) values ($1, $2, $3, $4) returning id", item.Name, item.Type, item.ActivityId, item.Priority).Scan(&tagId)
	if err != nil {
		return Tag{}, err
	}
	item.Id = tagId
	return item, nil
}
