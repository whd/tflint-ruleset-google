package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GoogleBigqueryDatasetNoProjectReadersRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "issue found",
			Content: `
resource "google_bigquery_dataset" "test" {
  dataset_id = "test"
  location   = "US"

  project = "project"
  access {
    role          = "READER"
    special_group = "projectReaders"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewGoogleBigqueryDatasetNoProjectReadersRule(),
					Message: `use of special group "projectReaders" is not allowed, use explicit GCPv2 workgroups instead`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 21},
						End:      hcl.Pos{Line: 9, Column: 37},
					},
				},
			},
		},
	}

	rule := NewGoogleBigqueryDatasetNoProjectReadersRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
