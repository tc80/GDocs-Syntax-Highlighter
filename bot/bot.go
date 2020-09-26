package main

import (
	"GDocs-Syntax-Highlighter/auth"
	"GDocs-Syntax-Highlighter/parser"
	"GDocs-Syntax-Highlighter/request"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func start(docID string, update time.Duration, verbose bool, docsService *docs.Service) {
	for {
		if verbose {
			log.Println("Fetching Google Document...")
		}
		doc, err := docsService.Documents.Get(docID).Do()
		if err != nil {
			log.Fatalf("Failed to get doc: %v", err)
		}

		var reqs []*docs.Request

		// process each instance of code found in the Google Doc
		instance := parser.GetCodeInstance(doc)

		if *instance.Shortcuts {
			// preprocess by replacing regex matches with specific strings
			for _, s := range instance.Lang.Shortcuts {
				reqs = append(reqs, instance.Replace(s)...)
			}
		}

		// attempt to format
		if instance.Format.Bold {
			// TODO: currently all directives are unbolded every time, so
			// there is no need to explicitly unbold at the moment
			//
			// unbold the #format directive to notify user that
			// the code was formatted or attempted to be formatted
			// reqs = append(reqs, request.SetBold(false, instance.Format.GetRange()))

			if instance.Lang.Format == nil {
				panic(fmt.Sprintf("no format func defined for language: `%s`", instance.Lang.Name))
			}
			if formatted, err := instance.Lang.Format(instance.Code); err != nil {
				// TODO: insert as Google Docs comment to notify of failure
				log.Printf("Failed to format: %v\n", err)
			} else {
				// After formatting, note that the new end index will be inaccurate
				// since the content length may have changed.
				// The end index will be updated later when we do further parsing.
				instance.Code = formatted

				// update for the new code string
				reqs = append(reqs, instance.UpdateCode()...)
			}
		}

		// attempt to run program
		if instance.Run.Bold {
			if instance.Lang.Run == nil {
				panic(fmt.Sprintf("no run func defined for language: `%s`", instance.Lang.Name))
			}
			res, err := instance.Lang.Run(instance.Code)
			if err != nil {
				// TODO: insert as Google Docs comment to notify of failure
				log.Printf("Failed to run: %v\n", err)
			} else {
				// TODO: insert run result as Google Docs comment
				fmt.Printf("RAN THE PROGRAM: %v\n", res)
			}
		}

		// map utf8 -> utf16, set end index
		instance.MapToUTF16()

		// set code foreground, code background, code font, code italics=false, doc background
		r, t := instance.GetRange(), instance.GetTheme()
		reqs = append(reqs, request.UpdateForegroundColor(t.CodeForeground, r))
		reqs = append(reqs, request.UpdateBackgroundColor(t.CodeBackground, r))
		reqs = append(reqs, request.UpdateHighlightColor(t.CodeHighlight, r))
		reqs = append(reqs, request.UpdateDocBackground(t.DocBackground))
		reqs = append(reqs, request.UpdateFont(*instance.Font, *instance.FontSize, r))
		reqs = append(reqs, request.SetItalics(false, r))

		for segmentID, seg := range instance.Segments {
			segRange := request.GetRange(seg.StartIndex, seg.EndIndex, segmentID)
			if seg.EndIndex == 1 {
				// empty header/footer (just `\n`), so replace config background color
				// with code's background to make the segment disappear
				reqs = append(reqs, request.UpdateBackgroundColor(t.CodeBackground, segRange))
				continue
			}
			reqs = append(reqs, request.UpdateForegroundColor(t.ConfigForeground, segRange))
			reqs = append(reqs, request.UpdateBackgroundColor(t.ConfigBackground, segRange))
			reqs = append(reqs, request.UpdateHighlightColor(t.ConfigHighlight, segRange))
			reqs = append(reqs, request.UpdateFont(t.ConfigFont, t.ConfigFontSize, segRange))
			reqs = append(reqs, request.SetItalics(t.ConfigItalics, segRange))
		}

		// remove ranges from instance.Code and add the requests to highlight them
		reqs = append(reqs, instance.RemoveRanges(t)...)

		// highlight code keywords using regexes
		for _, k := range t.Keywords {
			reqs = append(reqs, instance.Highlight(k.Regex, k.Color, "")...)
		}

		// update Google Document
		if len(reqs) > 0 {
			update := request.BatchUpdate(reqs)
			_, err := docsService.Documents.BatchUpdate(docID, update).Do()
			if err != nil {
				log.Printf("Failed to update Google Doc: %v\n", err)
			}
		}

		if verbose {
			log.Println("Sleeping...")
		}
		time.Sleep(update)
	}
}

func main() {
	log.Printf("Running...")

	var docID string
	var update int
	var verbose bool
	flag.StringVar(&docID, "doc", "", "Set the Google Document ID.")
	flag.IntVar(&update, "update", 1500, "Interval in milliseconds (>= 500) to update the Google Document.")
	flag.BoolVar(&verbose, "v", false, "Verbose mode.")
	flag.Parse()

	if docID == "" {
		flag.Usage()
		os.Exit(1)
	}

	if update < 500 {
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
	start(docID, time.Duration(update)*time.Millisecond, verbose, docsService)
}
