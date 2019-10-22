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

// A rune and its respective utf16 start and end indices
type char struct {
	startIndex int64 // the utf16 inclusive start index of the rune
	endIndex   int64 // the utf16 exclusive end index of the rune
	content    rune  // the rune
}

// A word as a string and its respective utf16 start and end indices
type word struct {
	startIndex int64  // the utf16 inclusive start index of the string
	endIndex   int64  // the utf16 exclusive end index of the string
	content    string // the string
}

func getWords(chars []char) []word {
	var words []word
	var b bytes.Buffer
	var startIndex, endIndex int64
	start := true
	for _, char := range chars {
		if unicode.IsSpace(char.content) {
			str := b.String()
			if len(str) > 0 {
				endIndex = char.startIndex // end index is start index of next
				words = append(words, word{startIndex, endIndex, str})
				start = true
				b.Reset()
			}
			continue
		}
		if start {
			startIndex = char.startIndex
			start = false
		}
		b.WriteRune(char.content)
	}
	str := b.String()
	if len(str) > 0 {
		endIndex = chars[len(chars)-1].endIndex // last char was not space, so last char's end index
		words = append(words, word{startIndex, endIndex, str})
	}
	return words
}

// Get the slice of chars, where each char holds a rune and its respective utf16 range
func getChars(doc *docs.Document) []char {
	var chars []char
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
						rUtf16 := utf16.Encode([]rune{r})                 // convert to utf16, since indices in GDocs API are utf16
						startIndex := index                               // start index of char
						index += int64(len(rUtf16))                       // add size of rune in utf16 format (now end index)
						chars = append(chars, char{startIndex, index, r}) // associate runes with ranges
					}
				}
			}
		}
	}
	return chars
}

func getColorRequest(r, g, b float64, startIndex, endIndex int64) *docs.Request {
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
							Red:   r,
							Blue:  b,
							Green: g,
						},
					},
				},
			},
		},
	}
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
		startIndex := chars[0].startIndex
		endIndex := chars[len(chars)-1].endIndex
		requests = append(requests, getFontRequest(startIndex, endIndex))

		words := getWords(chars)

		for _, w := range words {
			if strings.EqualFold(w.content, "public") {
				fmt.Println(w)
				requests = append(requests, getColorRequest(1, 0, 0, w.startIndex, w.endIndex))
			}
		}

		// for _, w := range words {
		// 	fmt.Printf("\nWord is (%v) (%v - %v)", w.content, w.startIndex, w.endIndex)
		// }

		update := getBatchUpdate(requests)
		response, err := docsService.Documents.BatchUpdate(docID, update).Do()
		_ = response

		// stop autocorrect?

		if err != nil {
			log.Fatalf("%v", err)
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
