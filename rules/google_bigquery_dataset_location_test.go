package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GoogleBigqueryDatasetLocation(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "basic",
			Content: `
resource "google_bigquery_dataset" "test" {
  location = "EU"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewGoogleBigqueryDatasetLocationRule(),
					Message: `expected location to be one of ["US"], got EU`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 18},
					},
				},
			},
		},
	}

	rule := NewGoogleBigqueryDatasetLocationRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
