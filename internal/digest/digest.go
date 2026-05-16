package digest

import (
	"context"
	"sort"

	"github.com/example/vvh-radar/internal/search"
	"github.com/example/vvh-radar/internal/storage"
)

type Item struct {
	storage.SearchResult
	Score int
}

func Build(ctx context.Context, store *storage.Store, topic, query string) ([]Item, error) {
	terms := search.Expand(query)
	rows, err := store.SearchAttachments(ctx, "", terms, topic)
	if err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(rows))
	for _, r := range rows {
		s := 0
		if r.AttachmentType.Valid {
			s += 2
		}
		if topic != "" && r.TopicTitle == topic {
			s += 2
		}
		for _, t := range terms {
			if t != "" && search.Normalize(r.Text) != "" {
				s++
				break
			}
		}
		items = append(items, Item{SearchResult: r, Score: s})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	if len(items) > 10 {
		items = items[:10]
	}
	return items, nil
}
