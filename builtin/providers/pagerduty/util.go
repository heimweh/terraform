package pagerduty

import (
	"fmt"
	"log"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

// Validate a value against a set of possible values
func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%#v is an invalid value for argument %s. Must be one of %#v", value, k, values))
		}
		return
	}
}

// getAllVendors retrieves all PagerDuty vendors and returns a list of []pagerduty.Vendor
func getAllVendors(client *pagerduty.Client) ([]pagerduty.Vendor, error) {

	var offset uint
	var totalCount int
	var vendors []pagerduty.Vendor
	var errorCount int

	for {
		o := &pagerduty.ListVendorOptions{
			APIListObject: pagerduty.APIListObject{
				Limit:  100,
				Total:  1,
				Offset: offset,
			},
		}

		resp, err := client.ListVendors(*o)

		if err != nil {
			errorCount++

			if errorCount < 4 {
				log.Printf("[WARN] Failed to retrieve vendors, retrying... (%d/%d)", errorCount, maxErrorCount)
				time.Sleep(1000 * time.Millisecond)
				continue
			}
			return nil, err
		}

		for _, v := range resp.Vendors {
			totalCount++
			vendors = append(vendors, v)
		}

		rOffset := uint(resp.Offset)
		returnedCount := uint(len(resp.Vendors))
		rTotal := uint(resp.Total)

		if resp.More && uint(totalCount) != uint(rTotal) {
			offset = returnedCount + rOffset
			continue
		}

		break
	}

	return vendors, nil
}
