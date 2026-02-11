package rules

import (
	"fmt"
	"path/filepath"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformPolicyDocLocationRule checks that aws_iam_policy_document data sources
// only appear in policies.tf.
type TerraformPolicyDocLocationRule struct {
	tflint.DefaultRule
}

// NewTerraformPolicyDocLocationRule creates a new rule.
func NewTerraformPolicyDocLocationRule() *TerraformPolicyDocLocationRule {
	return &TerraformPolicyDocLocationRule{}
}

// Name returns the rule name.
func (r *TerraformPolicyDocLocationRule) Name() string {
	return "terraform_policy_doc_location"
}

// Enabled returns whether the rule is enabled by default.
// Disabled by default as this is an organization-specific convention.
func (r *TerraformPolicyDocLocationRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *TerraformPolicyDocLocationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns a reference URL for the rule.
func (r *TerraformPolicyDocLocationRule) Link() string {
	return ""
}

// Check runs the rule against the given runner.
func (r *TerraformPolicyDocLocationRule) Check(runner tflint.Runner) error {
	schema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "data", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
		},
	}

	content, err := runner.GetModuleContent(schema, &tflint.GetModuleContentOption{})
	if err != nil {
		return fmt.Errorf("getting module content: %w", err)
	}

	for _, block := range content.Blocks {
		if len(block.Labels) < 2 || block.Labels[0] != "aws_iam_policy_document" {
			continue
		}

		basename := filepath.Base(block.DefRange.Filename)
		if basename != "policies.tf" {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("aws_iam_policy_document %q should be in policies.tf, found in %s", block.Labels[1], basename),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
