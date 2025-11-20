# Technology stack
This file describes the language and tools for the project.

## Language and Coding Conventions
- GO 1.25.0 for back end
- Use `snake_case` for file names.
- Front ends:
    - Angular CLI 20.2.0 (Node 22.18.0, npm 10.9.3)
    - VS-Code Extension in Typescript

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
│   ├── package.json
│   └── GEMINI.md 
│
├── GEMINI.md                     # Master context for LLM
├── .gitignore
├── README.md
└── note-app-vs-code/             # VS-Code Extension source code

## Tooling:
- git for version control 
- unit testing with Go standard library "testing"

## Commends:
- Go backend testing:
```bash
cd ./backend && go test ./...
```
- VS-Code frontend testing:
```bash
cd ./note-app-vs-code && npm test
```
