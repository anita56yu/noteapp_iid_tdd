# Technology stack
This file describes the language and tools for the project.

## Language and Coding Conventions
- GO 1.25.0 for back end
- Use `snake_case` for file names.
- Angular CLI 20.2.0 for front end
    - Node 22.18.0
    - Package Manager npm 10.9.3

## Directory structure
./
├── backend/                       # Go source code
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # Main application entry point
│   ├── internal/
│   │   ├── api/                   # HTTP handlers & routing
│   │   ├── domain/                # Core business models and logic
│   │   │   ├── note/              # Note aggregate
│   │   │   └── content/           # Content aggregate
│   │   ├── repository/            # Database interaction layer
│   │   │   ├── noterepo/
│   │   │   └── contentrepo/
│   │   └── usecase/               # Business logic orchestration
│   │       ├── noteuc/
│   │       └── contentuc/
│   ├── go.mod
│   └── go.sum
│
├── docs/                          # Contexts for LLM
│
├── frontend/                      # Angular source code
│   ├── src/
│   ├── angular.json
│   └── package.json
│
├── GEMINI.md                     # Master context for LLM
├── .gitignore
└── README.md

## Tooling:
- git for version control 
- unit testing with Go standard library "testing"

## Commends:
- Testing
```bash
cd ./backend && go test ./...
```
