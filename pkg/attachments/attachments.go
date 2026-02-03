package attachments

type Attachment struct {
	Id         int
	Name       string
	Link       string
	Priority   int
	ActivityId int
}

type AttachmentRepository interface {
	Get(activityId int) ([]Attachment, error)
	Add(item Attachment) (Attachment, error)
}
