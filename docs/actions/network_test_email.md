---
page_title: "network_test_email Action - das"
subcategory: "Organization"
description: |-
  Sends a test email of a purpose as the Network named by `network_id` (its branding).
---

# das_network_test_email (Action)

Sends a test email of a purpose as the Network named by `network_id` (its branding). Requires Terraform 1.14 or later.

API: `POST /api/networks/{networkId}/email/test-messages`.

## Example Usage

```terraform
action "das_network_test_email" "example" {
  config {
  network_id = "<networkId>"
  to         = "someone@example.com"
  }
}
```

Run it with `terraform apply -invoke=action.das_network_test_email.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.
- `to` (String) The recipient; any address (in test mode only verified addresses receive it).

### Optional

- `purpose` (String) Whose sender it goes out as: `platform.sign-in` (Sign-in), `platform.notifications` (Notifications). Default `platform.notifications`.

