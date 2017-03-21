package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyScheduleCreate,
		Read:   resourcePagerDutyScheduleRead,
		Update: resourcePagerDutyScheduleUpdate,
		Delete: resourcePagerDutyScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"layer": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotation_virtual_start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								d1, err := time.Parse(time.RFC3339, old)
								if err != nil {
									return false
								}

								d2, err := time.Parse(time.RFC3339, new)
								if err != nil {
									return false
								}

								return d1 == d2.Add(1*time.Hour)
							},
							StateFunc: func(v interface{}) string {
								switch v.(type) {
								case string:
									d, err := time.Parse(time.RFC3339, v.(string))
									if err != nil {
										return fmt.Sprintf("<failed>")
									}

									d.Add(-1 * time.Hour)

									return d.Format("2006-01-02T15:04:05.999999-07:00")
								default:
									return "<invalid>"
								}
							},
						},
						"rotation_turn_length_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"users": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"restriction": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time_of_day": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_day_of_week": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"duration_seconds": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildScheduleStruct(d *schema.ResourceData) *pagerduty.Schedule {
	scheduleLayers := d.Get("layer").([]interface{})

	schedule := pagerduty.Schedule{
		Name:           d.Get("name").(string),
		TimeZone:       d.Get("time_zone").(string),
		ScheduleLayers: expandScheduleLayers(scheduleLayers),
	}

	if attr, ok := d.GetOk("description"); ok {
		schedule.Description = attr.(string)
	}

	return &schedule
}

func resourcePagerDutyScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	schedule := buildScheduleStruct(d)

	log.Printf("[INFO] Creating PagerDuty schedule: %s", schedule.Name)

	schedule, err := client.CreateSchedule(*schedule)

	if err != nil {
		return err
	}

	d.SetId(schedule.ID)

	return resourcePagerDutyScheduleRead(d, meta)
}

func resourcePagerDutyScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty schedule: %s", d.Id())

	schedule, err := client.GetSchedule(d.Id(), pagerduty.GetScheduleOptions{})

	if err != nil {
		return err
	}

	d.Set("name", schedule.Name)
	d.Set("time_zone", schedule.TimeZone)
	d.Set("description", schedule.Description)

	if err := d.Set("layer", flattenScheduleLayers(schedule.ScheduleLayers)); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	schedule := buildScheduleStruct(d)

	log.Printf("[INFO] Updating PagerDuty schedule: %s", d.Id())

	if _, err := client.UpdateSchedule(d.Id(), *schedule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty schedule: %s", d.Id())

	if err := client.DeleteSchedule(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
