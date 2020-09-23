package parser

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"google.golang.org/api/docs/v1"
)

// Function to check if a particular
// rune is desired
type isRuneFunc func(r rune) bool

// Returns a function that will
// return true for any rune
func anyRune() isRuneFunc {
	return func(r1 rune) bool {
		return true
	}
}

// Returns a function that will
// return true only for the
// specified rune
func isRune(r1 rune) isRuneFunc {
	return func(r2 rune) bool {
		return r1 == r2
	}
}

// Gets the utf16 start and end indices of a target substring
// located in a utf8 string with a particular starting index offset.
func getUTF16SubstrIndices(target, utf8 string, offset int64) (startIndex, endIndex int64) {
	index := strings.Index(utf8, target)
	if index == -1 {
		panic(fmt.Sprintf("target `%s` not found in `%s`", target, utf8))
	}

	// add utf16 sizes until we reach the target's start
	startIndex += offset
	for _, r := range utf8[:index] {
		startIndex += GetUtf16RuneSize(r)
	}

	// endIndex is startIndex + utf16 size of target
	endIndex = startIndex
	for _, r := range target {
		endIndex += GetUtf16RuneSize(r)
	}

	return
}

// GetUtf16RuneSize gets the size of a rune in UTF-16 format
func GetUtf16RuneSize(r rune) int64 {
	rUtf16 := utf16.Encode([]rune{r}) // convert to utf16, since indices in GDocs API are utf16
	return int64(len(rUtf16))         // size of rune in utf16 format
}

// GetUtf16StringSize gets the size of a string in UTF-16 format
func GetUtf16StringSize(s string) int64 {
	var size int64
	for _, r := range s {
		size += GetUtf16RuneSize(r)
	}
	return size
}

// Gets all chars in a given document/line
// Gets the slice of all chars, where
// each Char holds a rune and its respective utf16 range
func getAllChars(doc *docs.Document) []*Char {
	var chars []*Char
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					index := par.StartIndex
					// iterate over runes
					for _, r := range par.TextRun.Content {
						size := GetUtf16RuneSize(r)                  // size of run in utf16 units
						chars = append(chars, &Char{index, size, r}) // associate runes with ranges
						index += size
					}
				}
			}
		}
	}
	return chars
}
