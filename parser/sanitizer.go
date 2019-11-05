package parser

import "google.golang.org/api/docs/v1"

// GetSanitizedWords gets the slice of Words
// defined between the document's defined
// START and (optionally) END region, ignoring
// comme

// need comments

// SanitizedDocument ..
type SanitizedDocument struct {
	Words    []Word
	Comments []Word
}

// linter - api ?
// gdocs comments to resolve

// https://github.com/collections/clean-code-linters

// should i try to use linters ?

//Sanitize ...
func Sanitize(doc *docs.Document) *SanitizedDocument {
	_ = getAllChars(doc)
	return nil
}
