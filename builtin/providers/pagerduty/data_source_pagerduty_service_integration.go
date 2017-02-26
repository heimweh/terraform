package pagerduty

import (
	"fmt"
	"log"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePagerDutyServiceIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyServiceIntegrationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyServiceIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty service integration")

	searchName := d.Get("name").(string)

	resp, err := client.ListServices(pagerduty.ListServiceOptions{})
	if err != nil {
		return err
	}

	var found *pagerduty.Integration

	for _, service := range resp.Services {
		for _, integration := range service.Integrations {
			if integration.Summary == searchName {
				found = &integration
				break
			}
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any service integration with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Summary)
	d.Set("integration_key", found.IntegrationKey)
	d.Set("integration_email", found.IntegrationEmail)

	return nil
}
