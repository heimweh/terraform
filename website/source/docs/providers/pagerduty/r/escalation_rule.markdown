---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_escalation_rule"
sidebar_current: "docs-pagerduty-resource-escalation_rule"
description: |-
  Creates and manages an escalation rule in PagerDuty.
---

# pagerduty\_escalation_rule

An [escalation rule](https://v2.developer.pagerduty.com/v2/page/api-reference#!/Escalation_Policies/get_escalation_policies_id_escalation_rules_escalation_rule_id) determines what user or schedule will be notified first, second, and so on when an incident is triggered.


## Example Usage

```
resource "pagerduty_user" "example" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
  teams = ["${pagerduty_team.example.id}"]
}

data "pagerduty_escalation_policy" "example" {
  name = "Engineering Escalation Policy"
}

resource "pagerduty_escalation_rule" "example" {
  escalation_policy_id = "${data.pagerduty_escalation_policy.example.id}"
  escalation_delay_in_minutes = 10

  target {
    id = "${pagerduty_user.example.id}"  
  }
}
```

## Argument Reference

The following arguments are supported:

* `escalation_policy_id` - (Required) The ID of the escalation policy the rule should belong to.
* `escalation_delay_in_minutes` - (Required) The number of minutes before an unacknowledged incident escalates away from this rule.
* `targets` - (Required) A target block. Target blocks documented below.


Targets (`target`) supports the following:

  * `type` - (Optional) Can be `user`, `schedule`, `user_reference` or `schedule_reference`. Defaults to `user_reference`
  * `id` - (Required) A target ID

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the escalation rule.

## Import

Escalation rules can be imported using the `id`, e.g.

```
$ terraform import pagerduty_escalation_rule.main PLBP09X
```
