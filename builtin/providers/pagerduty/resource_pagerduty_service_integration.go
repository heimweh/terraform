package pagerduty

import (
	"log"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePagerDutyServiceIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyServiceIntegrationCreate,
		Read:   resourcePagerDutyServiceIntegrationRead,
		Update: resourcePagerDutyServiceIntegrationUpdate,
		Delete: resourcePagerDutyServiceIntegrationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vendor"},
				ValidateFunc: validateValueFunc([]string{
					"aws_cloudwatch_inbound_integration",
					"cloudkick_inbound_integration",
					"event_transformer_api_inbound_integration",
					"generic_email_inbound_integration",
					"generic_events_api_inbound_integration",
					"keynote_inbound_integration",
					"nagios_inbound_integration",
					"pingdom_inbound_integration",
					"sql_monitor_inbound_integration",
				}),
			},
			"vendor": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"type"},
			},
			"integration_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"integration_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildServiceIntegrationStruct(d *schema.ResourceData) (*pagerduty.Integration, error) {
	name := d.Get("name").(string)
	service := d.Get("service").(string)

	integration := pagerduty.Integration{
		Name: name,
		APIObject: pagerduty.APIObject{
			Type: "service_integration",
		},
		Service: &pagerduty.APIObject{
			Type: "service",
			ID:   service,
		},
	}

	if attr, ok := d.GetOk("type"); ok {
		integration.Type = attr.(string)
	}

	if attr, ok := d.GetOk("vendor"); ok {
		integration.Vendor = &pagerduty.APIObject{
			ID:   attr.(string),
			Type: "vendor",
		}
	}

	return &integration, nil
}

func resourcePagerDutyServiceIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	integration, err := buildServiceIntegrationStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PagerDuty service integration %s", integration.Name)

	resp, err := client.CreateIntegration(integration.Service.ID, *integration)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePagerDutyServiceIntegrationRead(d, meta)
}

func resourcePagerDutyServiceIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty service integration %s", d.Id())

	service := d.Get("service").(string)

	o := &pagerduty.GetIntegrationOptions{}

	serviceIntegration, err := client.GetIntegration(service, d.Id(), *o)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", serviceIntegration.Name)
	d.Set("type", serviceIntegration.Type)
	d.Set("service", serviceIntegration.Service)
	d.Set("vendor", serviceIntegration.Vendor)
	d.Set("integration_key", serviceIntegration.IntegrationKey)
	d.Set("integration_email", serviceIntegration.IntegrationEmail)

	return nil
}

func resourcePagerDutyServiceIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	serviceIntegration, err := buildServiceIntegrationStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating PagerDuty service integration %s", d.Id())

	if _, err := client.UpdateIntegration(serviceIntegration.Service.ID, *serviceIntegration); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyServiceIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	service := d.Get("service").(string)

	log.Printf("[INFO] Removing PagerDuty service integration %s", d.Id())

	if err := client.DeleteIntegration(service, d.Id()); err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId("")

	return nil
}
