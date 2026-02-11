package rules

import (
	"strings"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformFileLineLimitRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string
		maxLines     int
		expectIssues int
	}{
		{
			name: "file under limit produces no issue",
			files: map[string]string{
				"main.tf": `resource "aws_instance" "example" {}`,
			},
			maxLines:     10,
			expectIssues: 0,
		},
		{
			name: "file at limit produces no issue",
			files: map[string]string{
				"main.tf": strings.Repeat("# comment\n", 5),
			},
			maxLines:     5,
			expectIssues: 0,
		},
		{
			name: "file over limit produces issue",
			files: map[string]string{
				"main.tf": strings.Repeat("# comment\n", 11),
			},
			maxLines:     10,
			expectIssues: 1,
		},
		{
			name: "empty file produces no issue",
			files: map[string]string{
				"main.tf": "",
			},
			maxLines:     10,
			expectIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, tt.files)
			rule := &TerraformFileLineLimitRule{MaxLines: tt.maxLines}

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := len(runner.Issues); got != tt.expectIssues {
				t.Errorf("expected %d issues, got %d: %v", tt.expectIssues, got, runner.Issues)
			}
		})
	}
}
