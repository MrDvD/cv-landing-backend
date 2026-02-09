package attachments

import (
	"cv-landing-backend/pkg/activity"
	"database/sql"
	"fmt"
	"strings"
)

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
			*attachment.Priority = int(raw.priority.Int16)
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

func (h AttachmentHandler) Edit(id int, ops []activity.EditField) (Attachment, error) {
	query, values := buildUpdateQuery(ops)
	rows, err := h.DB.Query(query, append(values, id)...)
	if err != nil {
		return Attachment{}, err
	}
	var attachment Attachment
	for rows.Next() {
		err := rows.Scan(&attachment.Id, &attachment.Name, &attachment.Link, &attachment.Priority, &attachment.ActivityId)
		if err != nil {
			return Attachment{}, err
		}
	}
	return attachment, nil
}

func buildUpdateQuery(ops []activity.EditField) (string, []any) {
	var query strings.Builder
	query.WriteString("update ATTACHMENTS")
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
	query.WriteString(" returning id, name, link, priority, activity_id")
	return query.String(), values
}
