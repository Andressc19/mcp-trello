# Contributing to mcp-trello

## Getting Started

1. Fork the repo
2. Clone your fork
3. Create a branch for your work

## Branch Strategy

We use a simplified GitHub Flow:

- `main` — stable, deployable code
- `feature/*` — new features
- `fix/*` — bug fixes

## Workflow

1. Create a branch from `main`:
   ```bash
   git checkout -b feature/my-new-feature
   # or
   git checkout -b fix/bug-description
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

3. Push and open a Pull Request:
   ```bash
   git push -u origin feature/my-new-feature
   ```

4. Wait for review, then merge

## Testing

Before submitting a PR, verify:

```bash
# Run tests
go test ./...

# Build
go build -o mcp-trello .

# Dry-run release (optional)
goreleaser release --dry-run --clean
```

## Releases

Tags with format `v*` automatically trigger releases via GitHub Actions:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This will:
- Build binaries for all platforms (Linux, macOS, Windows)
- Create a GitHub Release with artifacts
- Generate SHA256 checksums

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add tests for new functionality
