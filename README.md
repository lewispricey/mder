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
| `Ctrl+S` | Save (edit mode) |
| `Ctrl+E` | Toggle between view and edit mode |
| `q` | Quit (view mode) |
| `Ctrl+C` | Quit (press twice when unsaved changes exist) |

## Development

Prerequisites: Go 1.22+

```
go build ./...
go test ./... -race
```
