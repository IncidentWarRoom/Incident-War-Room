package incident

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID            uuid.UUID
	Title         string
	Severity      Severity
	Status        Status
	ChatID        int64
	TopicID       *int64 // linked Telegram Topic
	CreatedBy     *int64 // tg_user_id
	CreatedAt     time.Time
	ClosedAt      *time.Time
	TelegraphURLs []string // timeline Telegraph page URLs
	ReportURL     *string  // PDF report URL in object storage
}
