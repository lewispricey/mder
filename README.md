# mded — Terminal Markdown Editor

A terminal-based markdown viewer and live editor. Work in progress.

## Usage

```
mded [flags] <file>

Flags:
      --edit    Open in edit mode
      --view    Open in view-only mode (default)
```

## Key Bindings

| Key | Action |
|---|---|
| `q` | Quit (view mode) |
| `Ctrl+C` | Quit |

## Development

Prerequisites: Go 1.22+

```
go build ./...
go test ./... -race
```
