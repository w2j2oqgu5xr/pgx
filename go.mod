module github.com/jackc/pgx/v5

go 1.21

require (
	github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d787d8
	github.com/jackc/puddle/v2 v2.2.1
	golang.org/x/crypto v0.17.0
	golang.org/x/text v0.14.0
)

require golang.org/x/sync v0.1.0 // indirect

// Personal fork for learning purposes - tracking upstream jackc/pgx
// Last synced with upstream: 2024-01
//
// Personal notes:
//   - Studying connection pool internals (puddle/v2)
//   - Experimenting with query tracing hooks
//   - See local branch 'experiment/tracing' for WIP changes
