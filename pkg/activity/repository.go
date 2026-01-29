package activity

import (
	"database/sql"
)

type ActivityHandler struct {
	DB *sql.DB
}

func (h *ActivityHandler) GetAll() ([]Activity, error) {
	return h.getGeneric(nil)
}

func (h *ActivityHandler) GetAllOfType(activityType string) ([]Activity, error) {
	return h.getGeneric(&activityType)
}

func (h *ActivityHandler) getGeneric(rawType *string) ([]Activity, error) {
	activities := []Activity{}
	var rows *sql.Rows
	var err error
	if rawType == nil {
		rows, err = h.DB.Query("select id, name, subtitle, description, type, meta_label, date_start, date_end from ACTIVITIES")
	} else {
		rows, err = h.DB.Query("select id, name, subtitle, description, type, meta_label, date_start, date_end from ACTIVITIES where type = $1", *rawType)
	}
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
