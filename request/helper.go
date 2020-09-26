package request

import "strings"

// Joins the fields with commas.
func getFields(fields ...string) string {
	return strings.Join(fields, ",")
}
