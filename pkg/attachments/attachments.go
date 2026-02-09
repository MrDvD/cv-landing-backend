package attachments

import "cv-landing-backend/pkg/activity"

type Attachment struct {
	Id         int
	Name       string
	Link       string
	Priority   *int
	ActivityId int
}

type AttachmentRepository interface {
	Get(activityId int) ([]Attachment, error)
	Add(item Attachment) (Attachment, error)
	Remove(id int) error
	Edit(id int, ops []activity.EditField) (Attachment, error)
}
