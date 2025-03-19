# fantasySFC
Fantasy Football clone for GAA - written in Go.


## Project Structure
```
fantasySFC/
├── cmd/                  # Application entrypoints
│   ├── server/           # Main API server
│   └── scraper/          # Data collection tool
├── internal/             # Private application code
│   ├── core/             # Domain logic
│   │   ├── models/       # Domain entities (players, teams, etc.)
│   │   └── services/     # Business operations
│   ├── scraper/          # Web scraping logic
│   │   ├── client/       # HTTP client with retries
│   │   └── parser/       # HTML parsing
│   └── storage/          # Data persistence layer
│       ├── repository/   # Repository interfaces
│       ├── postgres/     # Relational storage (player stats, teams)
│       ├── mongodb/      # Document storage (game records, profiles)
│       └── redis/        # In-memory storage (leaderboards, caching)
├── pkg/                  # Public libraries
│   └── gaastats/         # Reusable GAA statistics
└── api/                  # API layer
    ├── rest/             # REST endpoints
    ├── grpc/             # gRPC service
    └── graphql/          # GraphQL resolvers
```
