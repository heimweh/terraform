package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourcePagerDutyService_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePagerDutyServiceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyService("pagerduty_service.test", "data.pagerduty_service.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyService(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a service ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the service %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyServiceConfig(rName string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
	name        = "TF user %[1]s"
	email       = "foo@example.com"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "test" {
	name        = "TF escalation policy %[1]s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = "${pagerduty_user.test.id}"
		}
	}
}

resource "pagerduty_service" "test" {
	name                    = "TF service %[1]s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = "${pagerduty_escalation_policy.test.id}"
	incident_urgency_rule {
		type    = "constant"
		urgency = "high"
	}
}

data "pagerduty_service" "by_name" {
  name = "${pagerduty_service.test.name}"
}
`, rName)
}
