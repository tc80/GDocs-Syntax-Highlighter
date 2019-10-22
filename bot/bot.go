package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"

	"GDocs-Syntax-Highlighter/auth"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const (
	// sometimes cannot find begin and end
	// need to fix
	beginSymbol        = "~~begin~~"
	endSymbol          = "~~end~~"
	courierNew         = "Courier New"
	foregroundColor    = "foregroundColor"
	weightedFontFamily = "weightedFontFamily"
)

var (
	black    = color{}
	red      = color{1, 0, 0}
	green    = color{0, 1, 0}
	blue     = color{0, 0, 1}
	keywords = map[string]color{
		"public": red,
		"static": blue,
		"void":   green,
	}
)

// A rune and its respective utf16 start and end indices
type char struct {
	index   int64 // the utf16 inclusive start index of the rune
	size    int64 // the size of the rune in utf16 units
	content rune  // the rune
}

// A word as a string and its respective utf16 start and end indices
type word struct {
	index   int64  // the utf16 inclusive start index of the string
	size    int64  // the size of the word in utf16 units
	content string // the string
}

// An RGB color
type color struct {
	red   float64 // the red value from 0.0 to 1.0
	green float64 // the green value from 0.0 to 1.0
	blue  float64 // the blue value from 0.0 to 1.0
}

func getWords(chars []*char) []*word {
	var words []*word
	var b bytes.Buffer
	var index int64
	start := true
	for _, char := range chars {
		if unicode.IsSpace(char.content) {
			str := b.String()
			if len(str) > 0 {
				size := getUtf16StringSize(str)
				words = append(words, &word{index, size, str})
				start = true
				b.Reset()
			}
			continue
		}
		if start {
			index = char.index
			start = false
		}
		b.WriteRune(char.content)
	}
	str := b.String()
	if len(str) > 0 {
		size := getUtf16StringSize(str)
		words = append(words, &word{index, size, str})
	}
	return words
}

// Get the size of a rune in UTF-16 format
func getUtf16RuneSize(r rune) int64 {
	rUtf16 := utf16.Encode([]rune{r}) // convert to utf16, since indices in GDocs API are utf16
	return int64(len(rUtf16))         // size of rune in utf16 format
}

// Get the size of a string in UTF-16 format
func getUtf16StringSize(s string) int64 {
	var size int64
	for _, r := range s {
		size += getUtf16RuneSize(r)
	}
	return size
}

// Get the slice of chars, where each char holds a rune and its respective utf16 range
func getChars(doc *docs.Document) []*char {
	var chars []*char
	begin := false
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					content := strings.TrimSpace(par.TextRun.Content)
					fmt.Println(content)
					if strings.EqualFold(content, endSymbol) {
						return chars
					}
					if !begin {
						if strings.EqualFold(content, beginSymbol) {
							begin = true
						}
						continue
					}
					index := par.StartIndex
					// iterate over runes
					for _, r := range par.TextRun.Content {
						size := getUtf16RuneSize(r)                  // size of run in utf16 units
						chars = append(chars, &char{index, size, r}) // associate runes with ranges
						index += size
					}
				}
			}
		}
	}
	return chars
}

// Gets a request to change the color of a range.
func getColorRequest(c color, startIndex, endIndex int64) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: foregroundColor,
			Range: &docs.Range{
				StartIndex: startIndex,
				EndIndex:   endIndex,
			},
			TextStyle: &docs.TextStyle{
				ForegroundColor: &docs.OptionalColor{
					Color: &docs.Color{
						RgbColor: &docs.RgbColor{
							Red:   c.red,
							Blue:  c.blue,
							Green: c.green,
						},
					},
				},
			},
		},
	}
}

// Gets the requests to delete a word and insert a new one in its place.
func getReplaceRequest(word *word, wordsAfter []*word, replace string) []*docs.Request {
	// delete word
	delete := &docs.Request{
		DeleteContentRange: &docs.DeleteContentRangeRequest{
			Range: &docs.Range{
				StartIndex: word.index,
				EndIndex:   word.index + word.size,
			},
		},
	}
	// insert replacement at deleted word's location
	insert := &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Text: replace,
			Location: &docs.Location{
				Index: word.index,
			},
		},
	}
	requests := []*docs.Request{delete, insert}
	newSize := getUtf16StringSize(replace)
	diff := newSize - word.size
	word.size = newSize
	// update ranges for words that follow this word
	for _, w := range wordsAfter {
		w.index += diff
	}
	return requests
}

func getFontRequest(startIndex, endIndex int64) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: weightedFontFamily,
			Range: &docs.Range{
				StartIndex: startIndex,
				EndIndex:   endIndex,
			},
			TextStyle: &docs.TextStyle{
				WeightedFontFamily: &docs.WeightedFontFamily{
					FontFamily: courierNew,
				},
			},
		},
	}
}

func getBatchUpdate(requests []*docs.Request) *docs.BatchUpdateDocumentRequest {
	return &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
}

func getRange(chars []*char) (int64, int64) {
	startIndex := chars[0].index
	lastChar := chars[len(chars)-1]
	endIndex := lastChar.index + lastChar.size
	return startIndex, endIndex
}

// for testing now
func start(docsService *docs.Service) {
	docID := "12Wqdvk_jk_pIfcN87o7X9EYvn4ukWRgNkpATpJwm1yM"
	for {
		fmt.Println("loop")
		doc, err := docsService.Documents.Get(docID).Do()
		if err != nil {
			log.Fatalf("Failed to get doc: %v", err)
		}

		chars := getChars(doc)

		if len(chars) == 0 {
			continue
		}

		var requests []*docs.Request
		startIndex, endIndex := getRange(chars)
		requests = append(requests, getColorRequest(black, startIndex, endIndex))
		requests = append(requests, getFontRequest(startIndex, endIndex))

		words := getWords(chars)
		for i, w := range words {
			lower := strings.ToLower(w.content)
			if c, ok := keywords[lower]; ok {
				if w.content != lower {
					// make lower
					requests = append(requests, getReplaceRequest(w, words[i+1:], lower)...)
				}
				requests = append(requests, getColorRequest(c, w.index, w.index+w.size))
				continue
			}
			if strings.EqualFold(lower, "psvm") {
				requests = append(requests, getReplaceRequest(w, words[i+1:], "public static void main(String[] args) {\n\n}")...)
			}
		}
		//requests = append(requests, getReplaceRequest(word{}, nil)...)
		// keywords
		// replaceall identifiers with a color?
		// make a set of identifiers
		// remove brackets from words, color brackets, then color words
		// formatting???
		// check word for identifier, not replaceall?
		// replace
		// if word is identifier, make lowercase/format????????

		// instead of replaceall, delete what is there and then insert at that location
		// will need to update indices of nearby elements?

		// for _, w := range words {
		// 	fmt.Printf("\nWord is (%v) (%v - %v)", w.content, w.startIndex, w.endIndex)
		// }

		update := getBatchUpdate(requests)
		response, err := docsService.Documents.BatchUpdate(docID, update).Do()
		_ = response

		// stop autocorrect?

		// HOW TO MAKE LOWERCASE???

		if err != nil {
			//log.Fatalf("%v", err)
		}
		time.Sleep(500 * time.Millisecond)
		//os.Exit(1)
	}
}

func main() {
	// get authorized client
	client, err := auth.GetAuthorizedClient()
	if err != nil {
		log.Fatalf("Failed to authorize client: %v", err)
	}

	// create docs service
	docsService, err := docs.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Docs service: %v", err)
	}

	// do stuff!
	start(docsService)
}
