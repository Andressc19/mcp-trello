# Contributing to mcp-trello

Thanks for contributing! mcp-trello enforces a strict **issue-first workflow** — every change starts with an approved issue.

---

## Contribution Workflow

```
Open Issue → Get status:approved → Open PR → Add type:* label → Review & Merge
```

### Step 1: Open an Issue

Use the correct template:
- **Bug Report** — for bugs
- **Feature Request** — for new features or improvements

> ⚠️ Blank issues are disabled. You must use a template.

Fill in all required fields. Your issue will automatically receive the `status:needs-review` label.

### Step 2: Wait for Approval

A maintainer will review the issue and add the `status:approved` label if it's accepted for implementation.

**Do not open a PR until the issue is approved.** Automated checks will block PRs that reference unapproved issues.

### Step 3: Open a Pull Request

Once the issue is approved:

1. Fork the repo and create a branch from `main`
2. Implement your change
3. Open a PR using the PR template — **link the approved issue** with `Closes #N`
4. Add exactly **one `type:*` label** to the PR (see label system below)

### Step 4: Automated PR Checks

Five checks run automatically on every PR:

#### PR Validation

| Check | What it verifies |
|-------|-----------------|
| **Check Issue Reference** | PR body contains `Closes #N`, `Fixes #N`, or `Resolves #N` |
| **Check Issue Has status:approved** | The linked issue has the `status:approved` label |
| **Check PR Has type:* Label** | PR has exactly one `type:*` label |

#### CI Tests

| Check | What it runs |
|-------|-------------|
| **Build** | `go build -o mcp-trello .` |
| **Unit Tests** | `go test ./...` |

All checks must pass before a PR can be merged.

---

## Label System

### Type Labels (required on every PR — pick exactly one)

| Label | Color | Use for |
|-------|-------|---------|
| `type:bug` | 🔴 | Bug fixes |
| `type:feature` | 🔵 | New features |
| `type:docs` | 🔵 | Documentation-only changes |
| `type:refactor` | 🟣 | Code refactoring with no behavior change |
| `type:chore` | ⚪ | Maintenance, tooling, dependencies |
| `type:breaking-change` | 🔴 | Breaking changes (requires major version bump) |

### Status Labels (set by maintainers)

| Label | Meaning |
|-------|---------|
| `status:needs-review` | Awaiting maintainer review (auto-applied to new issues) |
| `status:approved` | Approved for implementation — PRs can now be opened |
| `status:in-progress` | Actively being worked on — auto-exempt from stale bot |
| `status:blocked` | Blocked by another issue or external dependency |
| `status:stale` | No activity for 30 days — auto-applied by stale bot |
| `status:wontfix` | Intentionally not fixing — applied when closing stale/rejected items |

### Priority Labels (set by maintainers)

`priority:high`, `priority:medium`, `priority:low`

### Effort Labels (set by maintainers, for contributor guidance)

| Label | Meaning |
|-------|---------|
| `effort:small` | < 1 hour — good starting point for new contributors |
| `effort:medium` | 1–4 hours |
| `effort:large` | > 4 hours or spans multiple files |

---

## Branch Strategy

We use a simplified GitHub Flow:

- `main` — stable, deployable code
- `feature/*` — new features
- `fix/*` — bug fixes

---

## PR Rules

- Keep PR scope focused — one logical change per PR
- Use [conventional commits](https://www.conventionalcommits.org/) format
- Ensure all checks pass locally before pushing:
  - Build: `go build -o mcp-trello .`
  - Tests: `go test ./...`
- Update docs in the same PR when behavior changes
- Do not reference endpoints/scripts that do not exist in code
- Do not include `Co-Authored-By` trailers in commits

### Conventional Commit Format

```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

**Examples:**

```
feat(setup): add interactive TUI installer

fix(credentials): handle missing token gracefully

docs(readme): update installation instructions

refactor(trello): extract API client to separate package

chore(deps): bump github.com/charmbracelet/bubbletea to v0.26
```

Types map to labels: `feat` → `type:feature`, `fix` → `type:bug`, `docs` → `type:docs`, `refactor` → `type:refactor`, `chore` → `type:chore`.

---

## Testing

Before submitting a PR, verify:

```bash
# Run tests
go test ./...

# Build
go build -o mcp-trello .

# Format code
go fmt ./...

# Dry-run release (optional)
goreleaser release --dry-run --clean
```

---

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

---

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add tests for new functionality

---

## What Gets Closed Without Merging

- PRs opened without an approved issue
- PRs that fail CI and aren't updated within 30 days
- Issues that are vague, a duplicate, or belong in [Discussions](https://github.com/Andressc19/mcp-trello/discussions)
- Issues with no response to a maintainer question after 14 days
