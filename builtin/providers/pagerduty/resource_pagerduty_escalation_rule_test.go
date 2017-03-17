package pagerduty

import (
	"fmt"
	"log"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPagerDutyEscalationRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEscalationRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPagerDutyEscalationRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationRuleExists("pagerduty_escalation_policy.foo", "pagerduty_escalation_rule.foo"),
				),
			},
			resource.TestStep{
				Config: testAccCheckPagerDutyEscalationRuleConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationRuleExists("pagerduty_escalation_policy.foo", "pagerduty_escalation_rule.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEscalationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)

	var id string
	var escID string

	for _, r := range s.RootModule().Resources {
		switch r.Type {
		case "pagerduty_escalation_rule":
			id = r.Primary.ID
		case "pagerduty_escalation_policy":
			escID = r.Primary.ID
		default:
			continue
		}
	}

	log.Printf("[INFO] (DESTROY) Checking if %s/%s exists", escID, id)

	_, err := client.GetEscalationRule(escID, id, &pagerduty.GetEscalationRuleOptions{})

	if err == nil {
		return fmt.Errorf("Escalation Policy still exists")
	}

	return nil
}

func testAccCheckPagerDutyEscalationRuleExists(escID, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Escalation rule ID is set")
		}

		esc, ok := s.RootModule().Resources[escID]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if esc.Primary.ID == "" {
			return fmt.Errorf("No Escalation Policy ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		log.Printf("[INFO] (EXISTS) Checking if %s/%s exists", esc.Primary.ID, rs.Primary.ID)

		found, err := client.GetEscalationRule(esc.Primary.ID, rs.Primary.ID, &pagerduty.GetEscalationRuleOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Escalation rule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

const testAccCheckPagerDutyEscalationRuleConfig = `
resource "pagerduty_user" "foo" {
  name        = "foo"
  email       = "foo@bar.com"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_user" "bar" {
  name        = "foo"
  email       = "bar@foo.com"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "foo"
  description = "foo"
  num_loops   = 1

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = "${pagerduty_user.foo.id}"
    }
  }

	lifecycle {
	  ignore_changes = ["rule"]
	}
}

resource "pagerduty_escalation_rule" "foo" {
  escalation_policy_id = "${pagerduty_escalation_policy.foo.id}"
  escalation_delay_in_minutes = 30

	target {
	  type = "user_reference"
		id = "${pagerduty_user.bar.id}"
	}
}
`

const testAccCheckPagerDutyEscalationRuleConfigUpdated = `
resource "pagerduty_user" "foo" {
  name        = "foo"
  email       = "foo@bar.com"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_user" "bar" {
  name        = "foo"
  email       = "bar@foo.com"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "foo"
  description = "foo"
  num_loops   = 1

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = "${pagerduty_user.foo.id}"
    }
  }

	lifecycle {
	  ignore_changes = ["rule"]
	}
}

resource "pagerduty_escalation_rule" "foo" {
  escalation_policy_id = "${pagerduty_escalation_policy.foo.id}"
  escalation_delay_in_minutes = 60

	target {
	  type = "user_reference"
		id = "${pagerduty_user.foo.id}"
	}
}
`
