package pagerduty

import (
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePagerDutyEscalationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEscalationRuleCreate,
		Read:   resourcePagerDutyEscalationRuleRead,
		Update: resourcePagerDutyEscalationRuleUpdate,
		Delete: resourcePagerDutyEscalationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"escalation_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"escalation_delay_in_minutes": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"target": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user_reference",
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildEscalationRuleStruct(d *schema.ResourceData) *pagerduty.EscalationRule {
	return &pagerduty.EscalationRule{
		Delay:   uint(d.Get("escalation_delay_in_minutes").(int)),
		Targets: expandEscalationRuleTargets(d.Get("target").([]interface{})),
	}
}

func resourcePagerDutyEscalationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Creating PagerDuty escalation rule")

	escID := d.Get("escalation_policy_id").(string)

	rule := buildEscalationRuleStruct(d)

	escalationRule, err := client.CreateEscalationRule(escID, *rule)
	if err != nil {
		return err
	}

	d.SetId(escalationRule.ID)

	return resourcePagerDutyEscalationRuleRead(d, meta)
}

func resourcePagerDutyEscalationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty escalation rule: %s", d.Id())

	escID := d.Get("escalation_policy_id").(string)

	escalationRule, err := client.GetEscalationRule(escID, d.Id(), &pagerduty.GetEscalationRuleOptions{})
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("escalation_delay_in_minutes", escalationRule.Delay)

	return nil
}

func resourcePagerDutyEscalationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Updating PagerDuty escalation rule: %s", d.Id())

	escID := d.Get("escalation_policy_id").(string)

	rule := buildEscalationRuleStruct(d)

	if _, err := client.UpdateEscalationRule(escID, d.Id(), rule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyEscalationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty escalation rule: %s", d.Id())

	escID := d.Get("escalation_policy_id").(string)

	if err := client.DeleteEscalationRule(escID, d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
