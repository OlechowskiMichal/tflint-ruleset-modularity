# tflint-ruleset-modularity

A [TFLint](https://github.com/terraform-linters/tflint) ruleset that enforces Terraform module structure and modularity conventions.

## Rules

| Rule | Description | Default | Severity |
|------|-------------|---------|----------|
| `terraform_file_line_limit` | Limits `.tf` files to a maximum line count | Enabled (500 lines) | ERROR |
| `terraform_resource_file_limit` | Limits `resource`/`data` blocks per file | Enabled (5 blocks) | ERROR |
| `terraform_required_files` | Enforces required files exist in the module | Enabled (`variables.tf`, `outputs.tf`) | ERROR |
| `terraform_policy_doc_location` | Requires `aws_iam_policy_document` data sources in `policies.tf` | Disabled | ERROR |

### terraform_file_line_limit

Ensures no `.tf` file exceeds a maximum number of lines. Large files are harder to navigate and review.

```hcl
rule "terraform_file_line_limit" {
  enabled   = true
  max_lines = 500
}
```

### terraform_resource_file_limit

Ensures no single file contains too many `resource` and `data` blocks. Encourages splitting resources across files by concern.

```hcl
rule "terraform_resource_file_limit" {
  enabled       = true
  max_resources = 5
}
```

### terraform_required_files

Ensures required files exist in every Terraform module. Enforces a consistent module structure.

```hcl
rule "terraform_required_files" {
  enabled        = true
  required_files = ["variables.tf", "outputs.tf"]
}
```

### terraform_policy_doc_location

Enforces that `aws_iam_policy_document` data sources are defined in `policies.tf`. Disabled by default as this is an organization-specific convention.

```hcl
rule "terraform_policy_doc_location" {
  enabled = true
}
```

## Installation

Add the following to your `.tflint.hcl`:

```hcl
plugin "modularity" {
  enabled = true

  source  = "github.com/OlechowskiMichal/tflint-ruleset-modularity"
  version = "0.1.0"
}
```

Then run:

```bash
tflint --init
```

### Build from Source

```bash
task build:install
```

This builds the binary and installs it to `~/.tflint.d/plugins/`.

## Development

### Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev)

### Setup

```bash
task setup
```

### Commands

```bash
task test            # All tests with race detector
task test:unit       # Unit tests only
task lint            # Full lint
task lint:fast       # Lint only changed files
task build           # Build binary
task ci:verify       # Pre-commit checks (fmt + lint + test)
```

Run `task --list` for all available tasks.
