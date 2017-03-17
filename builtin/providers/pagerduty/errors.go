package pagerduty

import "strings"

func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "HTTP response code: 404")
}

func isUnauthorized(err error) bool {
	return strings.Contains(err.Error(), "HTTP response code: 401")
}

func isMissingAbility(err error) bool {
	return strings.Contains(err.Error(), "HTTP response code: 402")
}
