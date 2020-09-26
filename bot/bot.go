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
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func start(docID string, update time.Duration, verbose bool, docsService *docs.Service, driveService *drive.Service) {
	comments := drive.NewCommentsService(driveService)

	for {
		if verbose {
			log.Println("Fetching Google Document...")
		}
		doc, err := docsService.Documents.Get(docID).Do()
		if err != nil {
			log.Fatalf("Failed to get doc: %v", err)
		}

		var docsReqs []*docs.Request

		// process each instance of code found in the Google Doc
		instance := parser.GetCodeInstance(doc)

		if *instance.Shortcuts {
			// preprocess by replacing regex matches with specific strings
			for _, s := range instance.Lang.Shortcuts {
				docsReqs = append(docsReqs, instance.Replace(s)...)
			}
		}

		// attempt to format
		if instance.Format.Underlined {
			// un-underline the #format directive to notify user that
			// the code was formatted or attempted to be formatted
			docsReqs = append(docsReqs, request.SetUnderline(false, instance.Format.GetRange()))

			if instance.Lang.Format == nil {
				panic(fmt.Sprintf("no format func defined for language: `%s`", instance.Lang.Name))
			}
			if formatted, err := instance.Lang.Format(instance.Code); err != nil {
				log.Printf("Failed to format: %v\n", err)
				if _, err = request.CreateComment(fmt.Sprintf("Format Failure:\n%v", err), docID, comments).Do(); err != nil {
					log.Printf("Failed to create comment for format failure: %v\n", err)
				}
			} else {
				log.Println("Formatted the program.")

				// After formatting, note that the new end index will be inaccurate
				// since the content length may have changed.
				// The end index will be updated later when we do further parsing.
				instance.Code = formatted

				// update for the new code string
				docsReqs = append(docsReqs, instance.UpdateCode()...)
			}
		}

		// attempt to run program
		// in the future it might be good to run this on a separate thread,
		// but for now we will wait for it to complete
		if instance.Run.Underlined {
			// un-underline the #run directive to notify user that
			// the code was formatted or attempted to be formatted
			docsReqs = append(docsReqs, request.SetUnderline(false, instance.Run.GetRange()))

			if instance.Lang.Run == nil {
				panic(fmt.Sprintf("no run func defined for language: `%s`", instance.Lang.Name))
			}
			res, err := instance.Lang.Run(instance.Code)
			if err != nil {
				log.Printf("Failed to run: %v\n", err)
				if _, err = request.CreateComment(fmt.Sprintf("Run Internal Failure:\n%v", err), docID, comments).Do(); err != nil {
					log.Printf("Failed to create comment for run internal failure: %v\n", err)
				}
			} else {
				log.Printf("Ran the program (status=%d).\n", res.Status)
				if verbose {
					log.Printf("Program errors: %s\n", res.Errors)
					log.Printf("Program output: %s\n", res.Output)
				}
				if res.Errors == "" {
					if _, err = request.CreateComment(fmt.Sprintf("Run Success (status=%d):\n%s", res.Status, res.Output), docID, comments).Do(); err != nil {
						log.Printf("Failed to create comment for run success: %v\n", err)
					}
				} else {
					if _, err = request.CreateComment(fmt.Sprintf("Run Failure (status=%d):\n%s", res.Status, res.Errors), docID, comments).Do(); err != nil {
						log.Printf("Failed to create comment for run failure: %v\n", err)
					}
				}
			}
		}

		// map utf8 -> utf16, set end index
		instance.MapToUTF16()

		// set code foreground, code background, code font, code italics=false, doc background
		r, t := instance.GetRange(), instance.GetTheme()
		docsReqs = append(docsReqs, request.UpdateForegroundColor(t.CodeForeground, r))
		docsReqs = append(docsReqs, request.UpdateBackgroundColor(t.CodeBackground, r))
		docsReqs = append(docsReqs, request.UpdateHighlightColor(t.CodeHighlight, r))
		docsReqs = append(docsReqs, request.UpdateDocBackground(t.DocBackground))
		docsReqs = append(docsReqs, request.UpdateFont(*instance.Font, *instance.FontSize, r))
		docsReqs = append(docsReqs, request.ClearFormatting(r))

		for segmentID, seg := range instance.Segments {
			segRange := request.GetRange(seg.StartIndex, seg.EndIndex, segmentID)
			if seg.EndIndex == 1 {
				// empty header/footer (just `\n`), so replace config background color
				// with code's background to make the segment disappear
				docsReqs = append(docsReqs, request.UpdateBackgroundColor(t.CodeBackground, segRange))
				continue
			}
			docsReqs = append(docsReqs, request.UpdateForegroundColor(t.ConfigForeground, segRange))
			docsReqs = append(docsReqs, request.UpdateBackgroundColor(t.ConfigBackground, segRange))
			docsReqs = append(docsReqs, request.UpdateHighlightColor(t.ConfigHighlight, segRange))
			docsReqs = append(docsReqs, request.UpdateTextStyleExceptUnderline(
				t.ConfigFont, t.ConfigFontSize, t.ConfigItalics, t.ConfigBold, t.ConfigSmallCaps, t.ConfigStrikethrough, segRange,
			))
		}

		// remove ranges from instance.Code and add the requests to highlight them
		docsReqs = append(docsReqs, instance.RemoveRanges(t)...)

		// highlight code keywords using regexes
		for _, k := range t.Keywords {
			docsReqs = append(docsReqs, instance.Highlight(k.Regex, k.Color, "")...)
		}

		// update Google Document
		if len(docsReqs) > 0 {
			update := request.BatchUpdate(docsReqs)
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

	// create drive service
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Drive service: %v", err)
	}

	// start checking document
	start(docID, time.Duration(update)*time.Millisecond, verbose, docsService, driveService)
}
