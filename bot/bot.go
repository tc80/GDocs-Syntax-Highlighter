package main

import (
	"context"
	"fmt"
	"log"
	"unicode/utf16"

	"GDocs-Syntax-Highlighter/auth"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// A rune and its respective utf16 start and end indices
type char struct {
	startIndex int64 // the utf16 inclusive start index of the rune
	endIndex   int64 // the utf16 exclusive end index of the rune
	content    rune  // the rune
}

// Get the slice of chars, where each char holds a rune and its respective utf16 range
func getChars(doc *docs.Document) []char {
	var chars []char
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
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

// for testing now
func start(docsService *docs.Service) {
	docID := "12Wqdvk_jk_pIfcN87o7X9EYvn4ukWRgNkpATpJwm1yM"
	doc, err := docsService.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Failed to get doc: %v", err)
	}

	chars := getChars(doc)
	for _, c := range chars {
		fmt.Printf("\n%c (start: %v, end: %v)", c.content, c.startIndex, c.endIndex)
	}

	// with multiple updates what happens if you give the same range?
	// ex. remove (2 to 3), then remove (2 to 3)

	test := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{&docs.Request{
			InsertText: &docs.InsertTextRequest{
				Text: "a",
				Location: &docs.Location{
					Index: 1,
				},
			},
		}},
	}
	_ = test

	update := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{&docs.Request{
			UpdateTextStyle: &docs.UpdateTextStyleRequest{
				TextStyle: &docs.TextStyle{
					//Bold: true,
					ForegroundColor: &docs.OptionalColor{
						Color: &docs.Color{
							RgbColor: &docs.RgbColor{
								Red:   0.4,
								Green: 0.6,
								Blue:  0.6,
							},
						},
					},
				},
				Fields: "foregroundColor", // separate by commas
				Range: &docs.Range{ // need to keep track of ranges
					StartIndex: 3,
					EndIndex:   10,
				},
			},
		}},
	}
	_ = update
	response, err := docsService.Documents.BatchUpdate(docID, update).Do()
	_ = response

	// stop autocorrect?

	if err != nil {
		log.Fatalf("%v", err)
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
