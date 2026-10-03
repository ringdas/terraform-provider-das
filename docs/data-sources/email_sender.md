---
page_title: "email_sender Data Source - das"
subcategory: "Organization"
description: |-
  The sender of an email purpose (local part, display name, Reply-To) on the Organization's email domain. Read by its identity.
---

# das_email_sender (Data Source)

The sender of an email purpose (local part, display name, Reply-To) on the Organization's email domain. Read by its identity.

API: `GET /api/email/senders/{purpose}`.

## Example Usage

```terraform
data "das_email_sender" "example" {
  purpose = "example-purpose"
}
```

## Schema

### Required

- `purpose` (String) An email purpose, `<domain>.<kebab-name>`: `platform.sign-in` (Sign-in), `platform.notifications` (Notifications). A purpose the catalog does not know is 404 `NOT_FOUND`.

### Read-Only

- `address` (String) The From address, `<localPart>@<domain>`.
- `description` (String) Which messages it covers.
- `display_name` (String) The name receivers see.
- `is_default` (Boolean) True while the Organization has not set it: "DAS" <no-reply@<domain>>.
- `local_part` (String) The address before the @.
- `network_override` (Boolean) Whether a Network may set its own display name and Reply-To (setting `platform.email-sender-overrides`).
- `reply_to` (String) Where replies go; null: to the From address.
- `title` (String) What people call the purpose.

