package parser

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

// Sanitize ...
// func Sanitize(doc *docs.Document) *SanitizedDocument {
// 	chars := getAllChars(doc)
// }

// Gets all chars in a given document/line
// Gets the slice of all chars, where
// each Char holds a rune and its respective utf16 range
// func getAllChars(doc *docs.Document) []*Char {
// 	var chars []*Char
// 	for _, elem := range doc.Body.Content {
// 		if elem.Paragraph != nil {
// 			for _, par := range elem.Paragraph.Elements {
// 				if par.TextRun != nil {
// 					index := par.StartIndex
// 					// iterate over runes
// 					for _, r := range par.TextRun.Content {
// 						size := GetUtf16RuneSize(r)                  // size of run in utf16 units
// 						chars = append(chars, &Char{index, size, r}) // associate runes with ranges
// 						index += size
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return chars
// }
