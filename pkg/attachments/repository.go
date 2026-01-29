package attachments

import "database/sql"

type AttachmentHandler struct {
	DB *sql.DB
}

func (h *AttachmentHandler) Get(activityId int) ([]Attachment, error) {
	rows, err := h.DB.Query("select id, name, link, priority from ATTACHMENTS where activity_id = $1", activityId)
	if err != nil {
		return []Attachment{}, err
	}
	attachments := []Attachment{}
	for rows.Next() {
		raw := &struct {
			id       int
			name     string
			link     string
			priority sql.NullInt16
		}{}
		err := rows.Scan(&raw.id, &raw.name, &raw.link, &raw.priority)
		if err != nil {
			return []Attachment{}, err
		}
		attachment := Attachment{
			Id:         raw.id,
			Name:       raw.name,
			Link:       raw.link,
			ActivityId: activityId,
		}
		if raw.priority.Valid {
			attachment.Priority = int(raw.priority.Int16)
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}
