package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const defaultMaxResources = 5

// TerraformResourceFileLimitRule checks that no single file contains more than
// a maximum number of resource/data blocks.
type TerraformResourceFileLimitRule struct {
	tflint.DefaultRule

	MaxResources int
}

// NewTerraformResourceFileLimitRule creates a new rule with default configuration.
func NewTerraformResourceFileLimitRule() *TerraformResourceFileLimitRule {
	return &TerraformResourceFileLimitRule{MaxResources: defaultMaxResources}
}

// Name returns the rule name.
func (r *TerraformResourceFileLimitRule) Name() string {
	return "terraform_resource_file_limit"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformResourceFileLimitRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *TerraformResourceFileLimitRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns a reference URL for the rule.
func (r *TerraformResourceFileLimitRule) Link() string {
	return ""
}

type resourceFileLimitConfig struct {
	MaxResources int `hclext:"max_resources,optional"`
}

// Check runs the rule against the given runner.
func (r *TerraformResourceFileLimitRule) Check(runner tflint.Runner) error {
	config := &resourceFileLimitConfig{MaxResources: r.MaxResources}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return fmt.Errorf("decoding rule config: %w", err)
	}

	r.MaxResources = config.MaxResources

	schema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
			{Type: "data", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
		},
	}

	content, err := runner.GetModuleContent(schema, &tflint.GetModuleContentOption{})
	if err != nil {
		return fmt.Errorf("getting module content: %w", err)
	}

	// Deduplicate blocks by DefRange position. GetModuleContent expands
	// for_each/count at runtime, returning one block per instance — all
	// sharing the same DefRange from the original HCL definition. We
	// count unique definitions, not expanded instances.
	type defPos struct {
		file   string
		line   int
		column int
	}

	seen := make(map[defPos]struct{})
	fileCounts := make(map[string]int)
	fileFirstRange := make(map[string]hcl.Range)

	for _, block := range content.Blocks {
		pos := defPos{
			file:   block.DefRange.Filename,
			line:   block.DefRange.Start.Line,
			column: block.DefRange.Start.Column,
		}
		if _, exists := seen[pos]; exists {
			continue
		}
		seen[pos] = struct{}{}

		filename := block.DefRange.Filename

		fileCounts[filename]++
		if _, exists := fileFirstRange[filename]; !exists {
			fileFirstRange[filename] = block.DefRange
		}
	}

	for filename, count := range fileCounts {
		if count > r.MaxResources {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"%s has %d resource/data blocks, exceeding the limit of %d",
					filepath.Base(filename), count, r.MaxResources,
				),
				fileFirstRange[filename],
			); err != nil {
				return err
			}
		}
	}

	return nil
}
