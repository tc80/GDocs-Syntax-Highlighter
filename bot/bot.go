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

			fmt.Println(len(instance.Code))
			regex := regexp.MustCompile("\\bmap\\b")
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
		os.Exit(1)
		// os.Exit(1)

		// replace illegal character U+201C
		// replace illegal character U+201D

		// java, _ := style.GetLanguage("java")
		// chars, comms := parser.SeparateComments(java, chars)

		// if len(chars) == 0 {
		// 	continue
		// }

		// // underline if something is incorrect or spelled wrong ??/
		// // if something is wrong, highlight in red and for the error
		// // put in an actual gdocs comment
		// // when resolved, resolve the comment

		// var reqs []*docs.Request
		// startIndex, endIndex := parser.GetRange(chars)
		// reqs = append(reqs, requests.GetDocumentColorRequest(style.Black))
		// reqs = append(reqs, requests.GetBackColorRequest(style.Transparent, startIndex, endIndex))
		// reqs = append(reqs, requests.GetForeColorRequest(style.White, startIndex, endIndex))
		// reqs = append(reqs, requests.GetFontRequest(style.CourierNew, startIndex, endIndex))

		// for _, c := range comms {
		// 	fmt.Println(c)
		// 	reqs = append(reqs, requests.GetForeColorRequest(style.Green, c.Index, c.Index+c.Size))
		// }

		// words := parser.GetWords(chars)
		// for i, w := range words {
		// 	lower := strings.ToLower(w.Content)
		// 	if c, ok := java.Keywords[lower]; ok {
		// 		if w.Content != lower {
		// 			// make lower, probs should not iterate through words
		// 			// should split replace req into two methods maybe?
		// 			reqs = append(reqs, requests.GetReplaceRequest(w, words[i+1:], lower)...)
		// 		}
		// 		reqs = append(reqs, requests.GetForeColorRequest(c, w.Index, w.Index+w.Size))
		// 		continue
		// 	}
		// 	if replace, ok := java.Shortcuts[lower]; ok {
		// 		reqs = append(reqs, requests.GetReplaceRequest(w, words[i+1:], replace)...)
		// 	}
		// }
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
