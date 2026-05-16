package syncer

import (
	"context"
	"log/slog"

	"github.com/example/vvh-radar/internal/search"
	"github.com/example/vvh-radar/internal/storage"
	"github.com/example/vvh-radar/internal/vk"
)

type Syncer struct {
	store  *storage.Store
	client vk.Client
	logger *slog.Logger
}

func New(store *storage.Store, client vk.Client, logger *slog.Logger) *Syncer {
	return &Syncer{store: store, client: client, logger: logger}
}

func (s *Syncer) Run(ctx context.Context, groupID string) error {
	sourceID, err := s.store.EnsureSource(ctx, "fake-vk")
	if err != nil {
		return err
	}
	topics, err := s.client.GetTopics(ctx, groupID)
	if err != nil {
		return err
	}
	for _, t := range topics {
		t.SourceID = sourceID
		topicID, err := s.store.UpsertTopic(ctx, t)
		if err != nil {
			return err
		}
		msgs, atts, err := s.client.GetTopicMessages(ctx, groupID, t.ExternalID)
		if err != nil {
			return err
		}
		msgIDByExternal := map[int64]int64{}
		for _, m := range msgs {
			m.TopicID = topicID
			m.NormalizedText = search.Normalize(m.Text)
			mid, err := s.store.UpsertMessage(ctx, m)
			if err != nil {
				return err
			}
			msgIDByExternal[m.ExternalID] = mid
		}
		for _, a := range atts {
			if rid, ok := msgIDByExternal[a.MessageID]; ok {
				a.MessageID = rid
			}
			if err := s.store.InsertAttachment(ctx, a); err != nil {
				return err
			}
		}
	}
	s.logger.Info("sync completed")
	return nil
}
