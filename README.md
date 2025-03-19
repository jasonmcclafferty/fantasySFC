# fantasySFC
Fantasy Football clone for GAA - written in Go.


# Project Structure
fantasySFC/
├── cmd/                    # Application entrypoints
│   ├── server/             # The main server binary
│   ├── cli/                # CLI tools
│   └── scraper/            # Standalone scraper
├── internal/               # Private application code
│   ├── core/               # Domain models and business logic
│   │   ├── models/         # Domain entities
│   │   └── services/       # Business operations
│   ├── scraper/            # Web scraping logic
│   │   ├── client/         # HTTP client with retries
│   │   └── parser/         # HTML parsing
│   └── storage/            # Data persistence
│       ├── memory/         # In-memory implementation
│       └── postgres/       # Database implementation
├── pkg/                    # Public libraries
│   └── gaastats/           # Reusable GAA statistics
└── api/                    # Protocol interfaces
    ├── rest/               # REST API handlers
    ├── grpc/               # gRPC service definitions and handlers
    │   ├── proto/          # Protocol buffer definitions
    │   └── handlers/       # gRPC service implementations
    └── graphql/            # GraphQL schema and resolvers
        ├── schema/         # GraphQL schema definitions
        └── resolvers/      # GraphQL resolvers
