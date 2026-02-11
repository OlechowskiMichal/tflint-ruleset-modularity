# Agent Instructions

Instructions for AI coding agents working on the tflint-ruleset-modularity project.

## Project Overview

TFLint plugin ruleset enforcing Terraform module structure conventions. Built with the [TFLint plugin SDK](https://github.com/terraform-linters/tflint-plugin-sdk).

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

## Architecture

```
main.go              # Plugin entry point, registers rules with tflint.BuiltinRuleSet
rules/               # Rule implementations (one file per rule + test)
taskfiles/           # Modular task runner configs (build, test, lint, ci, release, setup)
lefthook/            # Git hook definitions
```

Each rule implements the `tflint.Rule` interface: `Name()`, `Enabled()`, `Severity()`, `Link()`, `Check(runner)`.

New rules: create `rules/terraform_<name>.go` + `rules/terraform_<name>_test.go`, register in `main.go`.

## Patterns to Follow

### Error Handling

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

### Rule Implementation

- Use `runner.GetFiles()` for file-level checks
- Use `runner.GetModuleContent(schema, opts)` for HCL block-level checks
- Emit issues via `runner.EmitIssue(rule, message, hcl.Range)`
- Use constants for default config values

### Testing Rules

- Use `helper.TestRunner(t, files)` to create test runners with in-memory files
- Table-driven tests with `t.Parallel()` at both top-level and subtest level
- Assert on `len(runner.Issues)` and issue messages

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
