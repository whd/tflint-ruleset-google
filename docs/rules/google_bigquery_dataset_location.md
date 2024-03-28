# google_bigquery_dataset_location

Disallow defining BQ datasets in any `location` other than the multi-region `US`.

## Example

```hcl
resource "google_bigquery_dataset" "foo" {
  location = "EU"
}
```

```
$ tflint

Error: expected location to be one of ["US"], got EU (google_bigquery_dataset_location)

  on main.tf line 2:
   2:   location = "EU"

```

## Why

Mozilla's primary data warehouse is defined in the `US` multi-region. Any BQ
resources defined elsewhere will not be query-able via the Data Platform and
its associated interfaces, such user-facing views.

## How To Fix

```hcl
resource "google_bigquery_dataset" "foo" {
  location = "US"
}
```
