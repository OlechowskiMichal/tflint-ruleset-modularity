// Package main provides the TFLint plugin entry point for the modularity ruleset.
package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/OlechowskiMichal/tflint-ruleset-modularity/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "modularity",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformFileLineLimitRule(),
				rules.NewTerraformResourceFileLimitRule(),
				rules.NewTerraformRequiredFilesRule(),
				rules.NewTerraformPolicyDocLocationRule(),
			},
		},
	})
}
