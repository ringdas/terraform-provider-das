---
page_title: "email Data Source - das"
subcategory: "Organization"
description: |-
  The Organization's email: domain, its checks and the sender of every purpose.
---

# das_email (Data Source)

The Organization's email: domain, its checks and the sender of every purpose.

API: `GET /api/email`.

## Example Usage

```terraform
data "das_email" "example" {
}
```

## Schema

### Read-Only

- `account` (Attributes) The sending account; null until a check reached the service.
- `attempted_at` (String) The last check, successful or not, ISO-8601 UTC.
- `checked_at` (String) The last check that reached the email service; the state below is from it, ISO-8601 UTC.
- `checks` (Attributes List) The checklist of the email domain, in order.
- `dns_records` (Attributes List) Every DNS record the domain needs, with what public DNS answers.
- `domain` (String) The Organization domain every sender's address is in.
- `error` (String) Why the last check failed; null when it succeeded.
- `mode` (String) The email service: Amazon SES, a mail directory (local development), or none (senders are stored, nothing is sent). One of `SES`, `DEV`, `NONE`.
- `ready` (Boolean) Whether DAS can send from the domain now: verified, signed with DKIM, out of test mode, not paused.
- `senders` (Attributes List) The sender of each email purpose.
- `sign_in_delivery` (Attributes) How the sign-in emails (invitations, codes) are sent now; null while no user pool is managed.

<a id="nestedatt--account"></a>
### Nested Schema for `account`

- `enforcement` (String) The provider's standing of the account: `HEALTHY`, `PROBATION` or `SHUTDOWN`.
- `max24_hour_send` (Number) Messages the account may send in 24 hours.
- `max_send_rate` (Number) Messages the account may send a second.
- `production_access` (Boolean) Out of test mode: mail reaches any address, not only verified ones.
- `sending_enabled` (Boolean) Whether the account sends at all.
- `sent_last24_hours` (Number) Messages sent in the last 24 hours.

<a id="nestedatt--checks"></a>
### Nested Schema for `checks`

- `id` (String) `IDENTITY` (the domain is verified), `DKIM` (messages are signed), `MAIL_FROM` (bounces return to the domain, SPF aligns), `DMARC`, `DNS_RECORDS` (every record resolves), `PRODUCTION_ACCESS` (mail reaches any address), `SENDING` (the account sends). One of `IDENTITY`, `DKIM`, `MAIL_FROM`, `DMARC`, `DNS_RECORDS`, `PRODUCTION_ACCESS`, `SENDING`.
- `status` (String) `PASS`, `WARN` (works, but act), `FAIL` (broken) or `UNKNOWN` (not checked). One of `PASS`, `WARN`, `FAIL`, `UNKNOWN`.
- `summary` (String) What was found, for people.

<a id="nestedatt--dns_records"></a>
### Nested Schema for `dns_records`

- `found` (List of String) What public DNS answered.
- `name` (String) The full name, without a trailing dot.
- `purpose` (String) What the record is for. One of `DKIM`, `MAIL_FROM`, `SPF`, `DMARC`.
- `status` (String) `FOUND` as it should be, `MISMATCH` something else, `MISSING` nothing, `NOT_CHECKED` (local development). One of `FOUND`, `MISMATCH`, `MISSING`, `NOT_CHECKED`.
- `type` (String) The record type. One of `CNAME`, `MX`, `TXT`.
- `value` (String) The value it must have; an MX value is `<priority> <host>`; for DMARC the recommended record (any policy of quarantine or reject passes).

<a id="nestedatt--senders"></a>
### Nested Schema for `senders`

- `address` (String) The From address, `<localPart>@<domain>`.
- `description` (String) Which messages it covers.
- `display_name` (String) The name receivers see.
- `is_default` (Boolean) True while the Organization has not set it: "DAS" <no-reply@<domain>>.
- `local_part` (String) The address before the @.
- `network_override` (Boolean) Whether a Network may set its own display name and Reply-To (setting `platform.email-sender-overrides`).
- `purpose` (String) The purpose: `platform.sign-in` (Sign-in), `platform.notifications` (Notifications).
- `reply_to` (String) Where replies go; null: to the From address.
- `title` (String) What people call the purpose.

<a id="nestedatt--sign_in_delivery"></a>
### Nested Schema for `sign_in_delivery`

- `applied_at` (String) When the pool last took it, ISO-8601 UTC.
- `attempted_at` (String) The duty's last attempt, ISO-8601 UTC.
- `error` (String) Why the pool refused the last attempt.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.
- `via` (String) `ORGANIZATION_DOMAIN`: from the sign-in sender, through the Organization domain; `BUILT_IN`: the user pool's own mailer (a low daily limit), until the domain is ready. One of `ORGANIZATION_DOMAIN`, `BUILT_IN`.

