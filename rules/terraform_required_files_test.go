package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformRequiredFilesRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string
		required     []string
		expectIssues int
	}{
		{
			name: "all required files present produces no issue",
			files: map[string]string{
				"main.tf":      `resource "aws_instance" "a" {}`,
				"variables.tf": `variable "name" {}`,
				"outputs.tf":   `output "id" { value = "" }`,
			},
			required:     []string{"variables.tf", "outputs.tf"},
			expectIssues: 0,
		},
		{
			name: "missing variables.tf produces one issue",
			files: map[string]string{
				"main.tf":    `resource "aws_instance" "a" {}`,
				"outputs.tf": `output "id" { value = "" }`,
			},
			required:     []string{"variables.tf", "outputs.tf"},
			expectIssues: 1,
		},
		{
			name: "missing outputs.tf produces one issue",
			files: map[string]string{
				"main.tf":      `resource "aws_instance" "a" {}`,
				"variables.tf": `variable "name" {}`,
			},
			required:     []string{"variables.tf", "outputs.tf"},
			expectIssues: 1,
		},
		{
			name: "missing both required files produces two issues",
			files: map[string]string{
				"main.tf": `resource "aws_instance" "a" {}`,
			},
			required:     []string{"variables.tf", "outputs.tf"},
			expectIssues: 2,
		},
		{
			name: "custom required files checked correctly",
			files: map[string]string{
				"main.tf": `resource "aws_instance" "a" {}`,
			},
			required:     []string{"providers.tf", "versions.tf"},
			expectIssues: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, tt.files)
			rule := &TerraformRequiredFilesRule{RequiredFiles: tt.required}

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := len(runner.Issues); got != tt.expectIssues {
				t.Errorf("expected %d issues, got %d: %v", tt.expectIssues, got, runner.Issues)
			}
		})
	}
}
