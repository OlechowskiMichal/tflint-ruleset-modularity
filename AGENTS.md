# Agent Instructions

Instructions for AI coding agents.

## Code Style

```yaml
Language: Go 1.25+
Errors: Wrap with context using fmt.Errorf("context: %w", err)
Naming: CamelCase (exported), camelCase (unexported)
```

## Getting Started

```bash
task setup            # First-time setup (install tools + git hooks)
task setup:deps-check # Verify all tools are installed
```

## Patterns to Follow

### Error Handling

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

## Testing

All commands use [Task](https://taskfile.dev). Run `task --list` to see all available tasks.

```bash
task test            # All tests with race detector
task test:unit       # Unit tests only
task test:race       # Race detector + short tests
task test:benchmark  # Benchmarks with memory stats
task lint            # Full lint
task lint:fast       # Lint only changed files
task lint:mocks      # Regenerate mocks (go generate)
task ci:verify       # Pre-commit checks (fmt + lint + test)
```

## Avoid

- Global state
- Panic for recoverable errors
- Naked returns
- Magic strings (use constants)
- Untested code paths

## Prefer

- Dependency injection
- Table-driven tests
- Interfaces for external dependencies
- Early returns
- Self-documenting code
