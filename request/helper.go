package request

import (
	"GDocs-Syntax-Highlighter/parser"

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

// Replace gets the requests to delete a Word and insert a new one in its place.
func Replace(word *parser.Word, wordsAfter []*parser.Word, replace string) []*docs.Request {
	// request to delete the Word
	delete := Delete(GetRange(word.Index, word.Index+word.Size))

	// request to insert the replacement at deleted Word's location
	insert := Insert(replace, word.Index)

	requests := []*docs.Request{delete, insert}
	newSize := parser.GetUtf16StringSize(replace)
	diff := newSize - word.Size
	word.Size = newSize
	// update ranges for Words that follow this Word
	for _, w := range wordsAfter {
		w.Index += diff
	}
	return requests
}

// BatchUpdate gets the batch request from a slice of requests.
func BatchUpdate(requests []*docs.Request) *docs.BatchUpdateDocumentRequest {
	return &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
}
