package pagerduty

import (
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourcePagerDutyUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyUserCreate,
		Read:   resourcePagerDutyUserRead,
		Update: resourcePagerDutyUserUpdate,
		Delete: resourcePagerDutyUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "user",
				ValidateFunc: validation.StringInSlice([]string{
					"admin",
					"limited_user",
					"owner",
					"read_only_user",
					"team_responder",
					"user",
				}, false),
			},

			"job_title": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			"teams": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"invitation_sent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceUser(d *schema.ResourceData) pagerduty.User {
	user := pagerduty.User{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
		APIObject: pagerduty.APIObject{
			ID: d.Id(),
		},
	}

	if v, ok := d.GetOk("color"); ok {
		user.Color = v.(string)
	}

	if v, ok := d.GetOk("role"); ok {
		role := v.(string)
		// Skip setting the role if the user is the owner of the account.
		// Can't change this through the API.
		if role != "owner" {
			user.Role = role
		}
	}

	if v, ok := d.GetOk("job_title"); ok {
		user.JobTitle = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		user.Description = v.(string)
	}

	return user
}

func resourcePagerDutyUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	user := resourceUser(d)

	log.Printf("[INFO] Creating PagerDuty user %s", user.Name)

	resp, err := client.CreateUser(user)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePagerDutyUserUpdate(d, meta)
}

func resourcePagerDutyUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty user %s", d.Id())

	resp, err := client.GetUser(d.Id(), pagerduty.GetUserOptions{})
	if err != nil {
		return handleNotFound(err, d, d.Id())
	}

	d.Set("avatar_url", resp.AvatarURL)
	d.Set("color", resp.Color)
	d.Set("description", resp.Description)
	d.Set("email", resp.Email)
	d.Set("html_url", resp.HTMLURL)
	d.Set("invitation_sent", resp.InvitationSent)
	d.Set("job_title", resp.JobTitle)
	d.Set("name", resp.Name)
	d.Set("role", resp.Role)
	d.Set("summary", resp.Summary)
	d.Set("teams", resp.Teams)
	d.Set("time_zone", resp.Timezone)

	contactMethods := map[string]interface{}{}
	for _, c := range resp.ContactMethods {
		contactMethods[c.Type] = c.Address
	}

	if len(contactMethods) > 0 {
		if err := d.Set("contact_methods", contactMethods); err != nil {
			return err
		}
	}

	return nil
}

func resourcePagerDutyUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	user := resourceUser(d)

	log.Printf("[INFO] Updating PagerDuty user %s", d.Id())

	if _, err := client.UpdateUser(user); err != nil {
		return err
	}

	if d.HasChange("teams") {
		o, n := d.GetChange("teams")

		if o == nil {
			o = new(schema.Set)
		}

		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		for _, team := range expandStringList(os.Difference(ns).List()) {
			if _, err := client.GetTeam(team); err != nil {
				log.Printf("[INFO] PagerDuty team: %s not found, removing dangling team reference for user %s", team, d.Id())
				continue
			}

			log.Printf("[INFO] Removing PagerDuty user %s from team: %s", d.Id(), team)

			if err := client.RemoveUserFromTeam(team, d.Id()); err != nil {
				return err
			}
		}

		for _, team := range expandStringList(ns.Difference(os).List()) {
			log.Printf("[INFO] Adding PagerDuty user %s to team: %s", d.Id(), team)

			if err := client.AddUserToTeam(team, d.Id()); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourcePagerDutyUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty user %s", d.Id())

	if err := client.DeleteUser(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
