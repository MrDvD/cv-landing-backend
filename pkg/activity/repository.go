package activity

import (
	"database/sql"
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
