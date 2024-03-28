# google_bigquery_dataset_no_project_readers

Disallow `projectReaders` special group usage in BQ `access` blocks.

## Example

```hcl
resource "google_bigquery_dataset" "test" {
  dataset_id = "test"
  location   = "US"

  project = "project"
  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }
  access {
    role          = "READER"
    special_group = "projectReaders"
  }
}
```

```
$ tflint

Error: use of special group "projectReaders" is not allowed, use explicit GCPv2 workgroups instead (google_bigquery_dataset_no_project_readers)

  on main.tf line 12:
  12:     special_group = "projectReaders"

```

## Why

The interaction between [Basic Roles](https://cloud.google.com/iam/docs/understanding-roles) and [BigQuery](https://cloud.google.com/bigquery/docs/access-control-basic-roles#project-basic-roles) makes it very easy to accidentally provision overbroad read access to (possibly workgroup-confidential) data. We generally want to provide engineering teams maximal visibility into GCP projects via the `Viewer` basic role but still being able to enforce workgroup-confidential dataset access where necessary.

See https://wiki.mozilla.org/Security/Data_Classification, https://mozilla-hub.atlassian.net/wiki/spaces/SRE/pages/27924789/Data+Access+Workgroups, and https://mozilla-hub.atlassian.net/browse/DSRE-1497 for additional background.

## How To Fix

You can explicitly grant the GCPv2 workgroups that are granted project-level `Viewer` access on a project `roles/bigquery.dataViewer` to a particular dataset via something like the following:

```
module "developers_workgroup" {
  source = "../../../../modules/workgroup"

  ids = ["workgroup:project/developers"]
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "test"
  location   = "US"

  project = "project"
  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }
  dynamic "access" {
    for_each = module.developers_workgroup.bigquery.read_acls

    content {
      role           = access.value.role
      user_by_email  = lookup(access.value, "user_by_email", null)
      group_by_email = lookup(access.value, "group_by_email", null)
      special_group  = lookup(access.value, "special_group", null)
    }
  }
}
```

See [here](https://github.com/mozilla-it/webservices-infra/blob/main/moztodon/tf/modules/moztodon_infra/bigquery.tf#L144) for a concrete example.

If absolutely sure the data is not sensitive (i.e. it is `workgroup:mozilla-confidential` or `PUBLIC`), or that the set of users that have basic `Viewer` on the project the dataset is defined in should have access to the data in the dataset, you can disable this rule via:

```
resource "google_bigquery_dataset" "test" {
  dataset_id = "test"
  location   = "US"

  project = "project"
  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }
  access {
    role = "READER"
    # tflint-ignore: google_bigquery_dataset_no_project_readers
    special_group = "projectReaders"
  }
}
```
