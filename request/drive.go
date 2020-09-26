package request

import "google.golang.org/api/drive/v3"

const (
	content = "content"
)

// CreateComment gets the *drive.CommentsCreateCall used to create
// a new Google Drive comment.
func CreateComment(comment, docID string, c *drive.CommentsService) *drive.CommentsCreateCall {
	return c.Create(docID, &drive.Comment{
		Content: comment,
	}).Fields(content)
}
