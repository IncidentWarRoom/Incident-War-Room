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
	TopicID       *int64 
	CreatedBy     *int64
	CreatedAt     time.Time
	ClosedAt      *time.Time
	TelegraphURLs []string
	ReportURL     *string  
}
