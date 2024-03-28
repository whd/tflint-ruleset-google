package rules

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TODO: Write the rule's description here
// GoogleBigqueryDatasetLocationRule checks ...
type GoogleBigqueryDatasetLocationRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewGoogleBigqueryDatasetLocationRule returns new rule with default attributes
func NewGoogleBigqueryDatasetLocationRule() *GoogleBigqueryDatasetLocationRule {
	return &GoogleBigqueryDatasetLocationRule{
		resourceType:  "google_bigquery_dataset",
		attributeName: "location",
	}
}

// Name returns the rule name
func (r *GoogleBigqueryDatasetLocationRule) Name() string {
	return "google_bigquery_dataset_location"
}

// Enabled returns whether the rule is enabled by default
func (r *GoogleBigqueryDatasetLocationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *GoogleBigqueryDatasetLocationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *GoogleBigqueryDatasetLocationRule) Link() string {
	return ""
}

func (r *GoogleBigqueryDatasetLocationRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(val string) error {
			validateFunc := validation.StringInSlice([]string{"US"}, false)

			_, errors := validateFunc(val, r.attributeName)
			for _, err := range errors {
				if err := runner.EmitIssue(r, err.Error(), attribute.Expr.Range()); err != nil {
					return err
				}
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
