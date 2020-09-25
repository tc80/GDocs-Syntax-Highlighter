package request

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

// GetRange gets a new *docs.Range for
// start and end indices.
func GetRange(start, end int64) *docs.Range {
	return &docs.Range{
		StartIndex: start,
		EndIndex:   end,
	}
}

// Joins the fields with commas.
func getFields(fields ...string) string {
	return strings.Join(fields, ",")
}

// BatchUpdate gets the batch request from a slice of requests.
func BatchUpdate(requests []*docs.Request) *docs.BatchUpdateDocumentRequest {
	return &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
}
