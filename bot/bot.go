package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf16"

	"GDocs-Syntax-Highlighter/auth"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const (
	beginSymbol        = "~~begin~~"
	endSymbol          = "~~end~~"
	courierNew         = "Courier New"
	weightedFontFamily = "weightedFontFamily"
)

// A rune and its respective utf16 start and end indices
type char struct {
	startIndex int64 // the utf16 inclusive start index of the rune
	endIndex   int64 // the utf16 exclusive end index of the rune
	content    rune  // the rune
}

type word struct {
	startIndex int64
	endIndex   int64
	content    string
}

func getWords(chars []char) []word {
	var words []word
	w := word{}
	fmt.Println(w.content)
	os.Exit(1)
	words = append(words, w)
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
		}}
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

		words := getWords(chars)

		var requests []*docs.Request
		startIndex := chars[0].startIndex
		endIndex := chars[len(chars)-1].endIndex
		requests = append(requests, getFontRequest(startIndex, endIndex))

		update := getBatchUpdate(requests)
		response, err := docsService.Documents.BatchUpdate(docID, update).Do()

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
