package vk

import (
	"context"
	"time"

	"github.com/example/vvh-radar/internal/app"
)

type Client interface {
	GetTopics(ctx context.Context, groupID string) ([]app.Topic, error)
	GetTopicMessages(ctx context.Context, groupID string, topicID int64) ([]app.Message, []app.Attachment, error)
}

type FakeClient struct{}

func (f *FakeClient) GetTopics(_ context.Context, groupID string) ([]app.Topic, error) {
	return []app.Topic{{ExternalID: 101, Title: "ОММ", URL: "https://vk.com/topic-" + groupID + "_101", RawJSON: `{"title":"ОММ"}`}}, nil
}

func (f *FakeClient) GetTopicMessages(_ context.Context, groupID string, topicID int64) ([]app.Message, []app.Attachment, error) {
	now := time.Now()
	msgs := []app.Message{
		{ExternalID: 1001, Author: "Иван", Text: "На зачет по ОММ будут билеты, смотрите файл программы.", URL: "https://vk.com/topic-" + groupID + "_" + "101?post=1001", PostedAt: now.Add(-72 * time.Hour), RawJSON: `{"id":1001}`},
		{ExternalID: 1002, Author: "Мария", Text: "Экзамен в пятницу, пересдача (хвост) через неделю.", URL: "https://vk.com/topic-" + groupID + "_" + "101?post=1002", PostedAt: now.Add(-48 * time.Hour), RawJSON: `{"id":1002}`},
		{ExternalID: 1003, Author: "Петр", Text: "Фото с билетами зачета прикрепил ниже.", URL: "https://vk.com/topic-" + groupID + "_" + "101?post=1003", PostedAt: now.Add(-24 * time.Hour), RawJSON: `{"id":1003}`},
	}
	atts := []app.Attachment{
		{MessageID: 1001, Type: "doc", Title: "Программа ОММ.pdf", URL: "https://example.local/omm_program.pdf", RawJSON: `{"type":"doc"}`},
		{MessageID: 1003, Type: "photo", Title: "Билеты зачет", URL: "https://example.local/omm_tickets.jpg", RawJSON: `{"type":"photo"}`},
	}
	_ = topicID
	return msgs, atts, nil
}
