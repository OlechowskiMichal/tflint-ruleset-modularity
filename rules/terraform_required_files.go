package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformRequiredFilesRule checks that required files exist in the module.
type TerraformRequiredFilesRule struct {
	tflint.DefaultRule

	RequiredFiles []string
}

// NewTerraformRequiredFilesRule creates a new rule with default configuration.
func NewTerraformRequiredFilesRule() *TerraformRequiredFilesRule {
	return &TerraformRequiredFilesRule{
		RequiredFiles: []string{"variables.tf", "outputs.tf"},
	}
}

// Name returns the rule name.
func (r *TerraformRequiredFilesRule) Name() string {
	return "terraform_required_files"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformRequiredFilesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformRequiredFilesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns a reference URL for the rule.
func (r *TerraformRequiredFilesRule) Link() string {
	return ""
}

type requiredFilesConfig struct {
	RequiredFiles []string `hclext:"required_files,optional"`
}

// Check runs the rule against the given runner.
func (r *TerraformRequiredFilesRule) Check(runner tflint.Runner) error {
	config := &requiredFilesConfig{RequiredFiles: r.RequiredFiles}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return fmt.Errorf("decoding rule config: %w", err)
	}

	r.RequiredFiles = config.RequiredFiles

	files, err := runner.GetFiles()
	if err != nil {
		return fmt.Errorf("getting files: %w", err)
	}

	existing := make(map[string]bool, len(files))
	for name := range files {
		existing[filepath.Base(name)] = true
	}

	// Use the first file as the issue location for missing files.
	var firstRange hcl.Range
	for name := range files {
		firstRange = hcl.Range{
			Filename: name,
			Start:    hcl.Pos{Line: 1, Column: 1},
			End:      hcl.Pos{Line: 1, Column: 1},
		}

		break
	}

	for _, required := range r.RequiredFiles {
		if !existing[required] {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("missing required file: %s", required),
				firstRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
