package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformResourceFileLimitRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string
		maxResources int
		expectIssues int
	}{
		{
			name: "single resource produces no issue",
			files: map[string]string{
				"main.tf": `resource "aws_instance" "a" {}`,
			},
			maxResources: 5,
			expectIssues: 0,
		},
		{
			name: "resources at limit produce no issue",
			files: map[string]string{
				"main.tf": `
resource "aws_instance" "a" {}
resource "aws_instance" "b" {}
`,
			},
			maxResources: 2,
			expectIssues: 0,
		},
		{
			name: "resources over limit produce issue",
			files: map[string]string{
				"main.tf": `
resource "aws_instance" "a" {}
resource "aws_instance" "b" {}
resource "aws_instance" "c" {}
`,
			},
			maxResources: 2,
			expectIssues: 1,
		},
		{
			name: "data sources count toward limit",
			files: map[string]string{
				"main.tf": `
resource "aws_instance" "a" {}
data "aws_ami" "b" {}
data "aws_vpc" "c" {}
`,
			},
			maxResources: 2,
			expectIssues: 1,
		},
		{
			name: "resources split across files stay under limit",
			files: map[string]string{
				"instances.tf": `
resource "aws_instance" "a" {}
resource "aws_instance" "b" {}
`,
				"networking.tf": `
resource "aws_vpc" "a" {}
resource "aws_subnet" "b" {}
`,
			},
			maxResources: 2,
			expectIssues: 0,
		},
		{
			name: "for_each resource counts as one block",
			files: map[string]string{
				"main.tf": `
resource "aws_budgets_budget" "account" {
  for_each = var.accounts
  name     = each.key
}
resource "aws_budgets_budget" "total" {
  name = "total"
}
`,
			},
			maxResources: 2,
			expectIssues: 0,
		},
		{
			name: "count resource counts as one block",
			files: map[string]string{
				"main.tf": `
resource "aws_instance" "a" {
  count = 3
}
resource "aws_instance" "b" {
  count = 5
}
resource "aws_instance" "c" {}
`,
			},
			maxResources: 3,
			expectIssues: 0,
		},
		{
			name: "for_each with data blocks counts correctly",
			files: map[string]string{
				"main.tf": `
data "aws_organizations_organization" "current" {}
data "aws_caller_identity" "current" {}
resource "aws_organizations_organizational_unit" "this" {
  for_each = var.ou_names
}
resource "aws_organizations_account" "this" {
  for_each = var.accounts
}
`,
			},
			maxResources: 4,
			expectIssues: 0,
		},
		{
			name: "no resources produces no issue",
			files: map[string]string{
				"main.tf": `variable "name" {}`,
			},
			maxResources: 5,
			expectIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, tt.files)
			rule := &TerraformResourceFileLimitRule{MaxResources: tt.maxResources}

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := len(runner.Issues); got != tt.expectIssues {
				t.Errorf("expected %d issues, got %d: %v", tt.expectIssues, got, runner.Issues)
			}
		})
	}
}
