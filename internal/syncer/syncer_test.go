package syncer

import (
	"context"
	"log/slog"
	"testing"

	"github.com/example/vvh-radar/internal/storage"
	"github.com/example/vvh-radar/internal/vk"
)

func TestFakeSync(t *testing.T) {
	s, _ := storage.Open(":memory:")
	defer s.Close()
	_ = s.ApplyMigrations(context.Background(), "../../migrations")
	sy := New(s, &vk.FakeClient{}, slog.Default())
	if err := sy.Run(context.Background(), "g"); err != nil {
		t.Fatal(err)
	}
	topics, _ := s.ListTopics(context.Background())
	if len(topics) == 0 {
		t.Fatal("expected topics")
	}
}
