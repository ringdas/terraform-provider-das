---
page_title: "test_email Action - das"
subcategory: "Organization"
description: |-
  Sends a test email of a purpose as the Organization.
---

# das_test_email (Action)

Sends a test email of a purpose as the Organization. Requires Terraform 1.14 or later.

API: `POST /api/email/test-messages`.

## Example Usage

```terraform
action "das_test_email" "example" {
  config {
  purpose = "platform.notifications"
  to      = "someone@example.com"
  }
}
```

Run it with `terraform apply -invoke=action.das_test_email.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `to` (String) The recipient; any address (in test mode only verified addresses receive it).

### Optional

- `purpose` (String) Whose sender it goes out as: `platform.sign-in` (Sign-in), `platform.notifications` (Notifications). Default `platform.notifications`.

