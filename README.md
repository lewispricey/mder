# mded — Terminal Markdown Editor

A terminal-based markdown viewer and live editor. Work in progress.

## Usage

```
mded [flags] <file>

Flags:
      --view    Open in view-only mode
      --edit    Open in edit mode (default)
```

## Development

Prerequisites: Go 1.22+

```
go build ./...
go test ./... -race
```
