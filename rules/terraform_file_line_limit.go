// Package rules implements TFLint rules for Terraform module structure enforcement.
package rules

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const defaultMaxLines = 500

// TerraformFileLineLimitRule checks that .tf files do not exceed a maximum line count.
type TerraformFileLineLimitRule struct {
	tflint.DefaultRule

	MaxLines int
}

// NewTerraformFileLineLimitRule creates a new rule with default configuration.
func NewTerraformFileLineLimitRule() *TerraformFileLineLimitRule {
	return &TerraformFileLineLimitRule{MaxLines: defaultMaxLines}
}

// Name returns the rule name.
func (r *TerraformFileLineLimitRule) Name() string {
	return "terraform_file_line_limit"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformFileLineLimitRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformFileLineLimitRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns a reference URL for the rule.
func (r *TerraformFileLineLimitRule) Link() string {
	return ""
}

// Check runs the rule against the given runner.
func (r *TerraformFileLineLimitRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return fmt.Errorf("getting files: %w", err)
	}

	for name, file := range files {
		if file == nil {
			continue
		}

		lines := bytes.Count(file.Bytes, []byte("\n"))
		if len(file.Bytes) > 0 && file.Bytes[len(file.Bytes)-1] != '\n' {
			lines++
		}

		if lines > r.MaxLines {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%s has %d lines, exceeding the limit of %d", filepath.Base(name), lines, r.MaxLines),
				hcl.Range{
					Filename: name,
					Start:    hcl.Pos{Line: 1, Column: 1},
					End:      hcl.Pos{Line: 1, Column: 1},
				},
			)
		}
	}

	return nil
}
