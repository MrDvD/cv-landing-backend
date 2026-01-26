package activity

import (
	"database/sql"
)

type RepositoryHandler struct {
	DB *sql.DB
}

func (h *RepositoryHandler) GetAll() ([]Activity, error) {
	return h.getGeneric(nil)
}

func (h *RepositoryHandler) GetAllOfType(activityType string) ([]Activity, error) {
	return h.getGeneric(&activityType)
}

func (h *RepositoryHandler) getGeneric(rawType *string) ([]Activity, error) {
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
			subtitle     string
			description  sql.NullString
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
			Name:      raw.name,
			Subtitle:  raw.subtitle,
			Type:      raw.activityType,
			DateStart: raw.dateStart,
		}
		if raw.description.Valid {
			activity.Description = &raw.description.String
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
