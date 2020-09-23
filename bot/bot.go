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
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const (
	sleepTime = time.Second * 10
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
			reqs = append(reqs, request.UpdateBackgroundColor(style.DarkThemeBackground, request.GetRange(instance.StartIndex, instance.EndIndex)))
			// // attempt to format if directive is bold and format is supported for this language
			// if instance.Format.Bold && instance.Lang.Format != nil {
			// 	if formatted, err := instance.Lang.Format(instance.Code); err != nil {
			// 		// TODO: insert as Google Docs comment?
			// 		log.Printf("Failed to format: %v\n", err)
			// 	} else {
			// 		// after formatting, need to know the new start and end indices???

			// 		reqs = append(reqs, requests.GetDeleteRequest(instance.StartIndex, instance.EndIndex))
			// 		reqs = append(reqs, requests.GetInsertRequest)
			// 		// 	req2 := requests.GetInsertRequest(res, begin)
			// 		instance.Code = formatted
			// 	}
			// }
			fmt.Println(instance.Lang.Format)
			_ = instance
			//request.GetDocumentColorRequest()
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

		// req := requests.GetForeColorRequest(style.Red, instances[0].Format.StartIndex, instances[0].Format.EndIndex)
		// update := requests.GetBatchUpdate([]*docs.Request{req})

		os.Exit(1)

		//res, err := style.FormatGo(text)
		//fmt.Println(res)
		//fmt.Println(err)

		// if err == nil {
		// 	req1 := requests.GetDeleteRequest(begin, end)
		// 	req2 := requests.GetInsertRequest(res, begin)
		// 	update := requests.GetBatchUpdate([]*docs.Request{req1, req2})
		// 	response, err := docsService.Documents.BatchUpdate(docID, update).Do()
		// 	_ = response
		// 	if err != nil {
		// 		fmt.Printf("\n\nERROR!!!!!: %v", err)
		// 		log.Fatalf("%v", err)
		// 	}
		// }

		// replace illegal character U+201C
		// replace illegal character U+201D

		// fmt.Println(text, begin, end, ok)
		// update := requests.GetBatchUpdate([]*docs.Request{requests.GetBackColorRequest(style.Red, begin, end)})
		//chars := parser.GetChars(doc)
		// response, err := docsService.Documents.BatchUpdate(docID, update).Do()
		// _ = response

		// if err != nil {
		// 	fmt.Printf("\n\nERROR!!!!!: %v", err)
		// 	log.Fatalf("%v", err)
		// }
		//os.Exit(1)

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
