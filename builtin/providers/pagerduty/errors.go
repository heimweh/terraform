package pagerduty

import "strings"

func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "Failed call API endpoint. HTTP response code: 404")
}

func isUnauthorized(err error) bool {
	return strings.Contains(err.Error(), "HTTP response code: 401")
}
