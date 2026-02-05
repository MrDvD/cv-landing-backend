package attachments

import "database/sql"

type AttachmentHandler struct {
	DB *sql.DB
}

func (h AttachmentHandler) Get(activityId int) ([]Attachment, error) {
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

func (h AttachmentHandler) Add(item Attachment) (Attachment, error) {
	var attachmentId int
	err := h.DB.QueryRow("insert into ATTACHMENTS(name, link, priority, activity_id) values ($1, $2, $3, $4) returning id", item.Name, item.Link, item.Priority, item.ActivityId).Scan(&attachmentId)
	if err != nil {
		return Attachment{}, err
	}
	item.Id = attachmentId
	return item, nil
}

func (h AttachmentHandler) Remove(id int) error {
	_, err := h.DB.Exec("delete from ATTACHMENTS where id = $1", id)
	return err
}
