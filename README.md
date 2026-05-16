# vvh-radar

vvh-radar is a local Go backend/CLI MVP that collects, stores, searches and exposes useful information from VK group discussion topics.

## Features (MVP)
- Fake VK client (no real API yet)
- SQLite storage with migrations and upserts
- Search with Russian normalization and synonym expansion
- Attachment filtering
- Digest scoring endpoint/command
- HTTP API + CLI

## Run locally
```bash
go mod tidy
go test ./...
go run ./cmd/vvh-radar init
go run ./cmd/vvh-radar sync
go run ./cmd/vvh-radar topics
go run ./cmd/vvh-radar search "зачет омм"
go run ./cmd/vvh-radar attachments --type photo --query "билеты зачет" --topic "ОММ"
go run ./cmd/vvh-radar digest --topic "ОММ" --query "зачет"
go run ./cmd/vvh-radar serve
```

## API
- `GET /health`
- `GET /api/topics`
- `POST /api/sync`
- `GET /api/search?q=...`
- `GET /api/attachments?type=photo&q=билеты&topic=ОММ`
- `GET /api/digest?topic=ОММ&q=зачет`

## Privacy
Do not commit tokens, secrets, local DB files, downloaded attachments, or private VK data.
