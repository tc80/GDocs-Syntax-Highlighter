package request

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

const (
	startIndex = "StartIndex"
	endIndex   = "EndIndex"
)

// GetRange gets a new *docs.Range for
// start and end indices.
func GetRange(start, end int64, segmentID string) *docs.Range {
	return &docs.Range{
		StartIndex: start,
		EndIndex:   end,
		SegmentId:  segmentID,
		// force send since a value of 0 in a header/footer
		// will be omitted in the JSON, causing a bad request
		ForceSendFields: []string{startIndex, endIndex},
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
