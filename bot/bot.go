package main

import (
	"GDocs-Syntax-Highlighter/auth"
	"GDocs-Syntax-Highlighter/parser"
	"GDocs-Syntax-Highlighter/requests"
	"GDocs-Syntax-Highlighter/style"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func start(docsService *docs.Service) {
	docID := "12Wqdvk_jk_pIfcN87o7X9EYvn4ukWRgNkpATpJwm1yM"
	for {
		fmt.Println("loop")
		doc, err := docsService.Documents.Get(docID).Do()
		if err != nil {
			log.Fatalf("Failed to get doc: %v", err)
		}

		chars := parser.GetChars(doc)

		java := style.GetLanguage("java")
		chars, comms := parser.SeparateComments(java, chars)

		if len(chars) == 0 {
			continue
		}

		var reqs []*docs.Request
		startIndex, endIndex := parser.GetRange(chars)
		reqs = append(reqs, requests.GetColorRequest(style.White, startIndex, endIndex))
		reqs = append(reqs, requests.GetFontRequest(style.CourierNew, startIndex, endIndex))

		for _, c := range comms {
			fmt.Println(c)
			reqs = append(reqs, requests.GetColorRequest(style.Green, c.Index, c.Index+c.Size))
		}

		words := parser.GetWords(chars)
		for i, w := range words {
			lower := strings.ToLower(w.Content)
			if c, ok := java.Keywords[lower]; ok {
				if w.Content != lower {
					// make lower, probs should not iterate through words
					// should split replace req into two methods maybe?
					reqs = append(reqs, requests.GetReplaceRequest(w, words[i+1:], lower)...)
				}
				reqs = append(reqs, requests.GetColorRequest(c, w.Index, w.Index+w.Size))
				continue
			}
			if replace, ok := java.Shortcuts[lower]; ok {
				reqs = append(reqs, requests.GetReplaceRequest(w, words[i+1:], replace)...)
			}
		}
		// if change is occuring before something, maybe wait to change its colors until change is
		// occuring after something so that its colors dont glitch

		//reqs = append(reqs, getReplaceRequest(Word{}, nil)...)
		// keywords
		// replaceall identifiers with a color?
		// make a set of identifiers
		// remove brackets from words, color brackets, then color words
		// formatting???
		// check Word for identifier, not replaceall?
		// replace
		// if Word is identifier, make lowercase/format????????

		// instead of replaceall, delete what is there and then insert at that location
		// will need to update indices of nearby elements?

		// for _, w := range words {
		// 	fmt.Printf("\nWord is (%v) (%v - %v)", w.content, w.startIndex, w.endIndex)
		// }

		//reqs = append(reqs, getDocumentRequest())
		update := requests.GetBatchUpdate(reqs)
		response, err := docsService.Documents.BatchUpdate(docID, update).Do()
		_ = response

		// stop autocorrect?

		// HOW TO MAKE LOWERCASE???

		if err != nil {
			fmt.Printf("\n\nERROR!!!!!: %v", err)
			log.Fatalf("%v", err)
		}
		time.Sleep(2000 * time.Millisecond)
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

	// start checking document
	start(docsService)
}
