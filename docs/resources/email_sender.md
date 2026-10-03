---
page_title: "email_sender Resource - das"
subcategory: "Organization"
description: |-
  The sender of an email purpose (local part, display name, Reply-To) on the Organization's email domain. Destroy goes back to the purpose's default sender.
---

# das_email_sender (Resource)

The sender of an email purpose (local part, display name, Reply-To) on the Organization's email domain. Destroy goes back to the purpose's default sender.

API: `PUT /api/email/senders/{purpose}` creates, `GET /api/email/senders/{purpose}` reads, `PUT /api/email/senders/{purpose}` updates, `DELETE /api/email/senders/{purpose}` deletes.

## Example Usage

```terraform
resource "das_email_sender" "example" {
  display_name = "DAS Notifications"
  local_part   = "notifications"
  purpose      = "platform.notifications"
}
```

## Schema

### Required

- `display_name` (String) The name receivers see; one line, no angle brackets, quotes or backslashes.
- `local_part` (String) The address before the @, stored in lower case, with no two dots in a row; the domain is always the Organization's.
- `purpose` (String) An email purpose, `<domain>.<kebab-name>`: `platform.sign-in` (Sign-in), `platform.notifications` (Notifications). A purpose the catalog does not know is 404 `NOT_FOUND`. Changing it replaces the resource.

### Optional

- `reply_to` (String) Where replies go; null: to the From address. Leave it out for null.

### Read-Only

- `address` (String) The From address, `<localPart>@<domain>`.
- `description` (String) Which messages it covers.
- `id` (String) The identity, purpose joined by `/` (the import id).
- `is_default` (Boolean) True while the Organization has not set it: "DAS" <no-reply@<domain>>.
- `network_override` (Boolean) Whether a Network may set its own display name and Reply-To (setting `platform.email-sender-overrides`).
- `title` (String) What people call the purpose.

## Import

```terraform
import {
  to = das_email_sender.example
  id = "<purpose>"
}
```
