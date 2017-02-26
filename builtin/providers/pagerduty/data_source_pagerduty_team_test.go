package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourcePagerDutyTeam_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePagerDutyTeamConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyTeam("pagerduty_team.test", "data.pagerduty_team.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyTeam(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a team ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the team %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyTeamConfig(rName string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "test" {
  name = "TF team %[1]s"
}

data "pagerduty_team" "by_name" {
  name = "${pagerduty_team.test.name}"
}
`, rName)
}
