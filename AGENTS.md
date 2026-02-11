# Agent Instructions

Instructions for AI coding agents.

## Code Style

```yaml
Language: Go 1.25+
Errors: Wrap with context using fmt.Errorf("context: %w", err)
Naming: CamelCase (exported), camelCase (unexported)
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
task test        # All tests with race detector
task test:unit   # Unit tests only
task lint        # Linters
task verify      # Pre-commit checks
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
