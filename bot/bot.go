package main

import (
	"GDocs-Syntax-Highlighter/auth"
	"GDocs-Syntax-Highlighter/parser"
	"GDocs-Syntax-Highlighter/request"
	"GDocs-Syntax-Highlighter/style"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const (
	sleepTime = time.Second * 1
)

func start(docID string, docsService *docs.Service) {
	for {
		doc, err := docsService.Documents.Get(docID).Do()
		if err != nil {
			log.Fatalf("Failed to get doc: %v", err)
		}

		var reqs []*docs.Request

		// process each instance of code found in the Google Doc
		for _, instance := range parser.GetCodeInstances(doc) {
			if *instance.Shortcuts {
				// note, need to update end index and make sure no shortcuts are in comments
				log.Println("TODO - process shortcuts")
			}

			lang := instance.Lang

			// attempt to format
			if instance.Format.Bold {
				// unbold the #format directive to notify user that
				// the code was formatted or attempted to be formatted
				reqs = append(reqs, request.SetBold(false, instance.Format.GetRange()))

				if lang.Format == nil {
					panic(fmt.Sprintf("no format func defined for language: `%s`", lang.Name))
				}
				if formatted, err := lang.Format(instance.Code); err != nil {
					// TODO: insert as Google Docs comment to notify of failure
					log.Printf("Failed to format: %v\n", err)
				} else {
					// delete the old text and insert the new text
					r := instance.GetRange()
					reqs = append(reqs, request.Delete(r))
					reqs = append(reqs, request.Insert(formatted, instance.StartIndex))

					// After formatting, note that the new end index will be inaccurate
					// since the content length may have changed.
					// The end index will be updated later when we do further parsing.
					instance.Code = formatted
				}
			}

			// ignore empty code
			if instance.Code == "" {
				continue
			}

			// map utf8 -> utf16, set end index
			instance.MapToUTF16()

			// set foreground, background, font, italics=false
			r, t := instance.GetRange(), instance.GetTheme()
			reqs = append(reqs, request.UpdateForegroundColor(t.Foreground, r))
			reqs = append(reqs, request.UpdateBackgroundColor(t.Background, r))
			reqs = append(reqs, request.UpdateFont(*instance.Font, *instance.FontSize, r))
			reqs = append(reqs, request.SetItalics(false, r))

			// remove ranges from instance.Code and add the requests to highlight them
			reqs = append(reqs, instance.RemoveRanges(t)...)

			regex := regexp.MustCompile("\\bint\\b")
			reqs = append(reqs, instance.Highlight(regex, style.DarkThemeBlue)...)

			// all shortcuts/formatting done, so now we just need to highlight
			// cannot be within comment
			// remove ranges, then do regex on remaining text, map index -> utf16 index
		}

		if len(reqs) > 0 {
			update := request.BatchUpdate(reqs)
			_, err := docsService.Documents.BatchUpdate(docID, update).Do()
			if err != nil {
				log.Fatalf("Failed to update Google Doc: %v\n", err)
			}
		}

		log.Println("Sleeping...")
		time.Sleep(sleepTime)

		// replace illegal character U+201C
		// replace illegal character U+201D

		os.Exit(1)
	}
}

func main() {
	var docID string
	flag.StringVar(&docID, "doc", "", "Set the Google Document ID.")
	flag.Parse()

	if docID == "" {
		flag.Usage()
		os.Exit(1)
	}

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
	start(docID, docsService)
}
