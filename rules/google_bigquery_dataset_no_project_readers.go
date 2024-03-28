// TODO new lint rule to require at least one access block with projectOwners
// also a dataset_access lint rule to preclude projectViewers usage or maybe
// specialGroup usage generally
package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// GoogleBigqueryDatasetNoProjectReadersRule checks whether ...
type GoogleBigqueryDatasetNoProjectReadersRule struct {
	tflint.DefaultRule
}

// NewGoogleBigqueryDatasetNoProjectReadersRule returns a new rule
func NewGoogleBigqueryDatasetNoProjectReadersRule() *GoogleBigqueryDatasetNoProjectReadersRule {
	return &GoogleBigqueryDatasetNoProjectReadersRule{}
}

// Name returns the rule name
func (r *GoogleBigqueryDatasetNoProjectReadersRule) Name() string {
	return "google_bigquery_dataset_no_project_readers"
}

// Enabled returns whether the rule is enabled by default
func (r *GoogleBigqueryDatasetNoProjectReadersRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *GoogleBigqueryDatasetNoProjectReadersRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *GoogleBigqueryDatasetNoProjectReadersRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *GoogleBigqueryDatasetNoProjectReadersRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("google_bigquery_dataset", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "access",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "role"},
						{Name: "special_group"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, rule := range resource.Body.Blocks {
			if attr, exists := rule.Body.Attributes["special_group"]; exists {
				err := runner.EvaluateExpr(attr.Expr, func(group string) error {
					if group == "projectReaders" {
						return runner.EmitIssue(
							r,
							fmt.Sprintf(`use of special group "%s" is not allowed, use explicit GCPv2 workgroups instead`, group),
							attr.Expr.Range(),
						)
					}
					return nil
				}, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
