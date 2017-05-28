package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func handleNotFound(err error, d *schema.ResourceData, resource string) error {
	if strings.Contains(err.Error(), "HTTP response code: 404") {
		// The resource doesn't exist anymore
		log.Printf("[WARN] Removing %s because it's gone", resource)
		d.SetId("")
		return nil
	}

	return fmt.Errorf("Error reading %s: %s", resource, err)
}

func isUnauthorized(err error) bool {
	return strings.Contains(err.Error(), "HTTP response code: 401")
}
