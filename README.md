# Database

A database engine written from scratch in Go.

This project implements a persistent, transactional database without leaning on
an existing storage or SQL library. It starts from a page-oriented file format
and builds upward: a copy-on-write B-tree, a key-value store on top of it, and
then a small relational layer with secondary indexes, transactions, and a query
language.

It exists as a systems project for learning Go in depth. Rather than toy
exercises, it deals with the things that make real software hard: on-disk
formats, crash recovery, concurrent readers and writers, and parsing.

## Features

### Storage engine

- Page-based file format with a master page and reference counting
- Copy-on-write B-tree with insert, delete, and range scan
- Durable persistence and atomic root-page updates
- Free list for reclaiming and reusing pages
- Node splitting, merging, and rebalancing

### Relational layer

- Typed rows and columns over the B-tree key space
- Secondary indexes maintained on write
- Range queries against primary and secondary indexes
- Atomic transactions with rollback on failure
- Concurrent readers and writers with locking
- A query language with a hand-written parser and executor

## Architecture

The engine is layered so each level depends only on the one below it.

```
      query language
    ┌─────────────────┐
    │  parser / exec  │
    ├─────────────────┤
    │  transactions   │
    ├─────────────────┤
    │ rows / indexes  │
    ├─────────────────┤
    │   key-value     │
    ├─────────────────┤
    │     B-tree      │
    ├─────────────────┤
    │ pages / free    │
    │      list       │
    ├─────────────────┤
    │   file / disk   │
    └─────────────────┘
```

## Project layout

```
.
├── cmd/
│   └── repl/       # interactive shell and entrypoint
├── internal/
│   ├── pages/      # page and file abstraction, free list
│   ├── btree/      # copy-on-write B-tree
│   ├── kv/         # key-value store over the B-tree
│   ├── table/      # rows, columns, secondary indexes
│   ├── txn/        # atomic transactions, locking
│   ├── parser/     # query language parser
│   └── exec/       # query execution
├── go.mod
└── README.md
```

Packages live under `internal/` because they are implementation details of the
engine, not a library API meant for other modules.

## Getting started

Requires Go 1.21 or newer.

```sh
git clone <repo-url> database
cd database

# download dependencies and build
go build ./...

# run the test suite
go test ./...

# run with the race detector
go test -race ./...

# start the interactive shell
go run ./cmd/repl
```

## Development

- Format with `gofmt`.
- Keep `go vet ./...` clean.
- New behavior is covered by tests in the same package.
- Performance-sensitive paths get benchmarks in the adjacent `*_test.go`.
- Errors carry context via `%w` wrapping and are inspected with `errors.Is`.

## References

- Build Your Own Database From Scratch in Go: https://build-your-own.org/database/
- Effective Go: https://go.dev/doc/effective_go
- Go Code Review Comments: https://go.dev/wiki/CodeReviewComments
- Database Internals, Alex Petrov
