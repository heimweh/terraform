package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePagerDutyVendor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyVendorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use field 'name' instead",
				ConflictsWith: []string{"name"},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyVendorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty vendor")

	var found *pagerduty.Vendor
	var serviceType string

	// Legacy lookup
	if _, ok := d.GetOk("name_regex"); ok {
		resp, err := getVendors(client)
		if err != nil {
			return err
		}

		r := regexp.MustCompile("(?i)" + d.Get("name_regex").(string))

		var vendors []pagerduty.Vendor
		var vendorNames []string

		for _, vendor := range resp {
			if r.MatchString(vendor.Name) {
				vendors = append(vendors, vendor)
				vendorNames = append(vendorNames, vendor.Name)
			}
		}

		if len(vendors) == 0 {
			return fmt.Errorf("Unable to locate any vendor using the regex string: %s", r.String())
		} else if len(vendors) > 1 {
			return fmt.Errorf("Your query returned more than one result using the regex string: %#v. Found vendors: %#v", r.String(), vendorNames)
		}

		found = &vendors[0]

		serviceType = found.GenericServiceType

	} else if _, ok := d.GetOk("name"); ok {
		searchName := d.Get("name").(string)

		o := &pagerduty.ListVendorOptions{
			Query: searchName,
		}

		resp, err := client.ListVendors(*o)
		if err != nil {
			return err
		}

		for _, vendor := range resp.Vendors {
			if strings.Contains(strings.ToLower(vendor.Name), strings.ToLower(searchName)) {
				found = &vendor
				break
			}
		}

		if found == nil {
			return fmt.Errorf("Unable to locate any vendor with the name: %s", searchName)
		}

		serviceType = found.GenericServiceType
	} else {
		return fmt.Errorf("Either name or name_regex must be specified!")
	}

	switch {
	case serviceType == "email":
		serviceType = "generic_email_inbound_integration"
	case serviceType == "api":
		serviceType = "generic_events_api_inbound_integration"
	default:
		return fmt.Errorf("Unknown service type")
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)
	d.Set("type", serviceType)

	return nil
}
