package cli

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vvh-radar/internal/api"
	"github.com/example/vvh-radar/internal/config"
	"github.com/example/vvh-radar/internal/digest"
	"github.com/example/vvh-radar/internal/search"
	"github.com/example/vvh-radar/internal/storage"
	"github.com/example/vvh-radar/internal/syncer"
	"github.com/example/vvh-radar/internal/vk"
)

func NewRoot() *cobra.Command {
	cfg := config.Load()
	root := &cobra.Command{Use: "vvh-radar"}
	open := func() *storage.Store {
		s, err := storage.Open(cfg.DBPath)
		if err != nil {
			panic(err)
		}
		return s
	}

	root.AddCommand(&cobra.Command{Use: "init", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		return s.ApplyMigrations(cmd.Context(), "migrations")
	}})
	root.AddCommand(&cobra.Command{Use: "sync", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		sy := syncer.New(s, &vk.FakeClient{}, slog.Default())
		return sy.Run(cmd.Context(), cfg.GroupID)
	}})
	root.AddCommand(&cobra.Command{Use: "topics", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		rows, _ := s.ListTopics(cmd.Context())
		for _, t := range rows {
			fmt.Println(t.Title, t.URL)
		}
		return nil
	}})
	root.AddCommand(&cobra.Command{Use: "search [query]", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		rows, _ := s.SearchMessages(cmd.Context(), search.Expand(args[0]))
		for _, r := range rows {
			fmt.Printf("[%s] %s\n%s\n%s\n", r.TopicTitle, r.PostedAt, r.Text, r.URL)
			if r.AttachmentURL.Valid {
				fmt.Printf("attachment(%s): %s\n", r.AttachmentType.String, r.AttachmentURL.String)
			}
			fmt.Println("---")
		}
		return nil
	}})

	var atType, atQuery, atTopic string
	att := &cobra.Command{Use: "attachments", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		rows, _ := s.SearchAttachments(cmd.Context(), atType, search.Expand(atQuery), atTopic)
		for _, r := range rows {
			fmt.Printf("%s | %s | %s\n", r.AttachmentType.String, r.AttachmentURL.String, r.Text)
		}
		return nil
	}}
	att.Flags().StringVar(&atType, "type", "", "")
	att.Flags().StringVar(&atQuery, "query", "", "")
	att.Flags().StringVar(&atTopic, "topic", "", "")
	root.AddCommand(att)

	var dTopic, dQuery string
	dc := &cobra.Command{Use: "digest", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		defer s.Close()
		items, _ := digest.Build(cmd.Context(), s, dTopic, dQuery)
		for _, it := range items {
			fmt.Printf("score=%d %s %s\n", it.Score, it.URL, it.Text)
		}
		return nil
	}}
	dc.Flags().StringVar(&dTopic, "topic", "", "")
	dc.Flags().StringVar(&dQuery, "query", "", "")
	root.AddCommand(dc)

	root.AddCommand(&cobra.Command{Use: "serve", RunE: func(cmd *cobra.Command, args []string) error {
		s := open()
		sy := syncer.New(s, &vk.FakeClient{}, slog.Default())
		return http.ListenAndServe(cfg.Addr, api.Router(s, sy, cfg.GroupID))
	}})
	root.SetContext(context.Background())
	_ = os.Setenv("TZ", "UTC")
	return root
}
