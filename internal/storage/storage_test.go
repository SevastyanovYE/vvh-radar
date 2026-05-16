package storage

import (
	"context"
	"testing"
	"time"

	"github.com/example/vvh-radar/internal/app"
)

func TestInsertAndSearch(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()
	_ = s.ApplyMigrations(context.Background(), "../../migrations")
	sid, _ := s.EnsureSource(context.Background(), "fake")
	tid, _ := s.UpsertTopic(context.Background(), app.Topic{SourceID: sid, ExternalID: 1, Title: "ОММ", URL: "u"})
	mid, _ := s.UpsertMessage(context.Background(), app.Message{TopicID: tid, ExternalID: 2, Text: "зачет скоро", NormalizedText: "зачет скоро", URL: "m", PostedAt: time.Now()})
	_ = s.InsertAttachment(context.Background(), app.Attachment{MessageID: mid, Type: "photo", URL: "a"})
	rows, _ := s.SearchMessages(context.Background(), []string{"зачет"})
	if len(rows) == 0 {
		t.Fatal("no rows")
	}
}
