package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/example/vvh-radar/internal/app"
)

type Store struct{ DB *sql.DB }

type SearchResult struct {
	TopicTitle     string
	PostedAt       string
	Text           string
	URL            string
	AttachmentType sql.NullString
	AttachmentURL  sql.NullString
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Close() error { return s.DB.Close() }

func (s *Store) ApplyMigrations(ctx context.Context, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err = s.DB.ExecContext(ctx, string(b)); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}
	return nil
}

func (s *Store) EnsureSource(ctx context.Context, name string) (int64, error) {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO sources(name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.DB.QueryRowContext(ctx, `SELECT id FROM sources WHERE name=?`, name).Scan(&id)
	return id, err
}

func (s *Store) UpsertTopic(ctx context.Context, t app.Topic) (int64, error) {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO topics(source_id, external_id, title, url, raw_json) VALUES(?,?,?,?,?)
ON CONFLICT(source_id, external_id) DO UPDATE SET title=excluded.title,url=excluded.url,raw_json=excluded.raw_json,updated_at=CURRENT_TIMESTAMP`, t.SourceID, t.ExternalID, t.Title, t.URL, t.RawJSON)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.DB.QueryRowContext(ctx, `SELECT id FROM topics WHERE source_id=? AND external_id=?`, t.SourceID, t.ExternalID).Scan(&id)
	return id, err
}

func (s *Store) UpsertMessage(ctx context.Context, m app.Message) (int64, error) {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO messages(topic_id, external_id, author, text, normalized_text, url, posted_at, raw_json)
VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(topic_id, external_id) DO UPDATE SET author=excluded.author,text=excluded.text,normalized_text=excluded.normalized_text,url=excluded.url,posted_at=excluded.posted_at,raw_json=excluded.raw_json,updated_at=CURRENT_TIMESTAMP`,
		m.TopicID, m.ExternalID, m.Author, m.Text, m.NormalizedText, m.URL, m.PostedAt, m.RawJSON)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.DB.QueryRowContext(ctx, `SELECT id FROM messages WHERE topic_id=? AND external_id=?`, m.TopicID, m.ExternalID).Scan(&id)
	return id, err
}

func (s *Store) InsertAttachment(ctx context.Context, a app.Attachment) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO attachments(message_id,type,title,url,raw_json) VALUES(?,?,?,?,?)
ON CONFLICT(message_id,type,url) DO UPDATE SET title=excluded.title,raw_json=excluded.raw_json`, a.MessageID, a.Type, a.Title, a.URL, a.RawJSON)
	return err
}

func (s *Store) ListTopics(ctx context.Context) ([]app.Topic, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,title,url FROM topics ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []app.Topic{}
	for rows.Next() {
		var t app.Topic
		_ = rows.Scan(&t.ID, &t.Title, &t.URL)
		out = append(out, t)
	}
	return out, nil
}

func (s *Store) SearchMessages(ctx context.Context, terms []string) ([]SearchResult, error) {
	if len(terms) == 0 {
		return nil, nil
	}
	clauses := make([]string, 0, len(terms))
	args := make([]any, 0, len(terms))
	for _, t := range terms {
		clauses = append(clauses, "m.normalized_text LIKE ?")
		args = append(args, "%"+t+"%")
	}
	q := `SELECT t.title, m.posted_at, m.text, m.url, a.type, a.url FROM messages m JOIN topics t ON t.id=m.topic_id LEFT JOIN attachments a ON a.message_id=m.id WHERE ` + strings.Join(clauses, " OR ") + ` ORDER BY m.posted_at DESC`
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SearchResult{}
	for rows.Next() {
		var r SearchResult
		_ = rows.Scan(&r.TopicTitle, &r.PostedAt, &r.Text, &r.URL, &r.AttachmentType, &r.AttachmentURL)
		out = append(out, r)
	}
	return out, nil
}

func (s *Store) SearchAttachments(ctx context.Context, typ string, terms []string, topic string) ([]SearchResult, error) {
	q := `SELECT t.title,m.posted_at,m.text,m.url,a.type,a.url FROM attachments a JOIN messages m ON m.id=a.message_id JOIN topics t ON t.id=m.topic_id WHERE 1=1`
	args := []any{}
	if typ != "" {
		q += ` AND a.type=?`
		args = append(args, typ)
	}
	if topic != "" {
		q += ` AND t.title=?`
		args = append(args, topic)
	}
	for _, t := range terms {
		q += ` AND m.normalized_text LIKE ?`
		args = append(args, "%"+t+"%")
	}
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SearchResult{}
	for rows.Next() {
		var r SearchResult
		_ = rows.Scan(&r.TopicTitle, &r.PostedAt, &r.Text, &r.URL, &r.AttachmentType, &r.AttachmentURL)
		out = append(out, r)
	}
	return out, nil
}
