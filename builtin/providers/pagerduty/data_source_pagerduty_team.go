package pagerduty

import (
	"fmt"
	"log"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty team")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListTeamOptions{
		Query: searchName,
	}

	resp, err := client.ListTeams(*o)
	if err != nil {
		return err
	}

	var found *pagerduty.Team

	for _, team := range resp.Teams {
		if team.Name == searchName {
			found = &team
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any team with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)

	return nil
}
