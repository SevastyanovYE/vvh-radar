package app

import "time"

type Source struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type Topic struct {
	ID         int64
	SourceID   int64
	ExternalID int64
	Title      string
	URL        string
	RawJSON    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Message struct {
	ID             int64
	TopicID        int64
	ExternalID     int64
	Author         string
	Text           string
	NormalizedText string
	URL            string
	PostedAt       time.Time
	RawJSON        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Attachment struct {
	ID        int64
	MessageID int64
	Type      string
	Title     string
	URL       string
	RawJSON   string
	CreatedAt time.Time
}

type SyncRun struct {
	ID         int64
	SourceID   int64
	Status     string
	StartedAt  time.Time
	FinishedAt *time.Time
	Error      string
}
