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
