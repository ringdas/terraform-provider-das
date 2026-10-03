---
page_title: "email_domain_check Action - das"
subcategory: "Organization"
description: |-
  Checks the Organization's email domain now (DNS and the email service) and records the result.
---

# das_email_domain_check (Action)

Checks the Organization's email domain now (DNS and the email service) and records the result. Requires Terraform 1.14 or later.

API: `POST /api/email/checks`.

## Example Usage

```terraform
action "das_email_domain_check" "example" {
  config {
  }
}
```

Run it with `terraform apply -invoke=action.das_email_domain_check.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

