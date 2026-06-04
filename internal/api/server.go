package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/vvh-radar/internal/digest"
	"github.com/example/vvh-radar/internal/search"
	"github.com/example/vvh-radar/internal/storage"
	"github.com/example/vvh-radar/internal/syncer"
)

func Router(store *storage.Store, sync *syncer.Syncer, groupID string) http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/api/topics", func(w http.ResponseWriter, req *http.Request) {
		t, _ := store.ListTopics(req.Context())
		_ = json.NewEncoder(w).Encode(t)
	})
	r.Post("/api/sync", func(w http.ResponseWriter, req *http.Request) {
		err := sync.Run(req.Context(), groupID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	})
	r.Get("/api/search", func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		rows, _ := store.SearchMessages(req.Context(), search.Expand(q))
		_ = json.NewEncoder(w).Encode(rows)
	})
	r.Get("/api/attachments", func(w http.ResponseWriter, req *http.Request) {
		typ := req.URL.Query().Get("type")
		q := req.URL.Query().Get("q")
		topic := req.URL.Query().Get("topic")
		rows, _ := store.SearchAttachments(req.Context(), typ, search.Expand(q), topic)
		_ = json.NewEncoder(w).Encode(rows)
	})
	r.Get("/api/digest", func(w http.ResponseWriter, req *http.Request) {
		topic := req.URL.Query().Get("topic")
		q := req.URL.Query().Get("q")
		rows, _ := digest.Build(req.Context(), store, topic, q)
		_ = json.NewEncoder(w).Encode(rows)
	})
	return r
}
