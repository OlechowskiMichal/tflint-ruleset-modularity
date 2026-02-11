package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformPolicyDocLocationRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string
		expectIssues int
	}{
		{
			name: "policy doc in policies.tf produces no issue",
			files: map[string]string{
				"policies.tf": `data "aws_iam_policy_document" "example" {}`,
			},
			expectIssues: 0,
		},
		{
			name: "policy doc in main.tf produces issue",
			files: map[string]string{
				"main.tf": `data "aws_iam_policy_document" "example" {}`,
			},
			expectIssues: 1,
		},
		{
			name: "non-policy data source in main.tf produces no issue",
			files: map[string]string{
				"main.tf": `data "aws_ami" "example" {}`,
			},
			expectIssues: 0,
		},
		{
			name: "multiple policy docs in wrong file produces multiple issues",
			files: map[string]string{
				"iam.tf": `
data "aws_iam_policy_document" "read" {}
data "aws_iam_policy_document" "write" {}
`,
			},
			expectIssues: 2,
		},
		{
			name: "mixed correct and incorrect locations",
			files: map[string]string{
				"policies.tf": `data "aws_iam_policy_document" "correct" {}`,
				"main.tf":     `data "aws_iam_policy_document" "wrong" {}`,
			},
			expectIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, tt.files)
			rule := NewTerraformPolicyDocLocationRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := len(runner.Issues); got != tt.expectIssues {
				t.Errorf("expected %d issues, got %d: %v", tt.expectIssues, got, runner.Issues)
			}
		})
	}
}
