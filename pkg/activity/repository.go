package activity

import (
	"database/sql"
	"fmt"
	"strings"
)

type ActivityHandler struct {
	DB *sql.DB
}

func (h ActivityHandler) GetAllOfType(activityType string) ([]Activity, error) {
	activities := []Activity{}
	rows, err := h.DB.Query("select id, name, subtitle, description, type, meta_label, date_start, date_end from ACTIVITIES where type = $1", activityType)
	if err != nil {
		return []Activity{}, err
	}
	for rows.Next() {
		raw := &struct {
			id           int
			name         string
			subtitle     sql.NullString
			description  string
			activityType string
			metaLabel    sql.NullString
			dateStart    string
			dateEnd      sql.NullString
		}{}
		err := rows.Scan(&raw.id, &raw.name, &raw.subtitle, &raw.description, &raw.activityType, &raw.metaLabel, &raw.dateStart, &raw.dateEnd)
		if err != nil {
			return []Activity{}, err
		}
		activity := Activity{
			Id:          raw.id,
			Name:        raw.name,
			Description: raw.description,
			Type:        raw.activityType,
			DateStart:   raw.dateStart,
		}
		if raw.subtitle.Valid {
			activity.Subtitle = &raw.subtitle.String
		}
		if raw.metaLabel.Valid {
			activity.MetaLabel = &raw.metaLabel.String
		}
		if raw.dateEnd.Valid {
			activity.DateEnd = &raw.dateEnd.String
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (h ActivityHandler) Add(item Activity) (Activity, error) {
	var activityId int
	err := h.DB.QueryRow("insert into ACTIVITIES(name, subtitle, description, type, meta_label, date_start, date_end) values ($1, $2, $3, $4, $5, $6, $7) returning id", item.Name, item.Subtitle, item.Description, item.Type, item.MetaLabel, item.DateStart, item.DateEnd).Scan(&activityId)
	if err != nil {
		return Activity{}, err
	}
	item.Id = activityId
	return item, nil
}

func (h ActivityHandler) Remove(id int) error {
	_, err := h.DB.Exec("delete from ACTIVITIES where id = $1", id)
	return err
}

func (h ActivityHandler) Edit(id int, ops []EditField) (Activity, error) {
	query, values := buildUpdateQuery(ops)
	rows, err := h.DB.Query(query, append(values, id)...)
	if err != nil {
		return Activity{}, err
	}
	var activity Activity
	for rows.Next() {
		err := rows.Scan(&activity.Id, &activity.Name, &activity.Subtitle, &activity.Description, &activity.Type, &activity.MetaLabel, &activity.DateStart, &activity.DateEnd)
		if err != nil {
			return Activity{}, err
		}
	}
	return activity, nil
}

func buildUpdateQuery(ops []EditField) (string, []any) {
	var query strings.Builder
	query.WriteString("update ACTIVITIES")
	setters := []string{}
	values := []any{}
	for i, op := range ops {
		setters = append(setters, fmt.Sprintf("%s = $%d", op.Name, i+1))
		values = append(values, op.Value)
	}
	if len(setters) != 0 {
		query.WriteString(" SET ")
	}
	query.WriteString(strings.Join(setters, ", "))
	query.WriteString(fmt.Sprintf(" where id = $%d", len(ops)+1))
	query.WriteString(" returning id, name, subtitle, description, type, meta_label, date_start, date_end")
	return query.String(), values
}
