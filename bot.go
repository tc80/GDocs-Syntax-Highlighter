package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const (
	scope           = docs.DriveScope    // needed for editing GDrive files
	stateToken      = "state-token"      // used for requesting a new token
	credentialsPath = "credentials.json" // client secret
	tokenPath       = "token.json"       // token path, needs to change if scope changes
)

// assuming endindex = startindex + 1
// what is rune is wider than 1 index?
type char struct {
	startIndex int64
	endIndex   int64
	content    rune
}

// Authorizes the client with an API token.
func authorizeClient(config *oauth2.Config) *http.Client {
	token, err := checkForToken()
	if err != nil {
		fmt.Println("Unable to locate local token, attempting to get token from web.")
		token = requestNewToken(config)
		cacheToken(token)
	}
	return config.Client(context.Background(), token)
}

// Request a new token from the Docs API.
func requestNewToken(config *oauth2.Config) *oauth2.Token {
	// get authorization code
	fmt.Printf("Enter auth code from: \n%v\n", config.AuthCodeURL(stateToken, oauth2.AccessTypeOffline))
	var auth string
	_, err := fmt.Scan(&auth)
	if err != nil {
		log.Fatalf("Failed to scan auth code: %v", err)
	}

	// get new token using auth code, passing empty context (same as TODO())
	token, err := config.Exchange(oauth2.NoContext, auth)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	return token
}

// Checks if the client already has a local token.
func checkForToken() (*oauth2.Token, error) {
	// open file for reading
	file, err := os.Open(tokenPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}

	// parse token json into Token
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

// Cache the new token.
func cacheToken(token *oauth2.Token) {
	// open file for writing, allow it to be read/written to, create if doesn't exist, truncate it to length 0
	file, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer file.Close()
	if err != nil {
		log.Fatalf("Failed to write token to file: %v", err)
	}

	// encode the token into json
	json.NewEncoder(file).Encode(token)
}

// for testing now
func subscribe(docsService *docs.Service) {
	docID := "12Wqdvk_jk_pIfcN87o7X9EYvn4ukWRgNkpATpJwm1yM"
	doc, err := docsService.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Failed to get doc: %v", err)
	}

	// replaceall Public -> public, etc If -> if depending on lang
	// var b bytes.Buffer
	// fmt.Printf("\n\nEND: %v", c)
	// b.WriteString(par.TextRun.Content)

	var chars []char
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					var index int64 = par.StartIndex
					// parStartIndex := par.StartIndex
					// var offset int64 = 0
					for _, r := range par.TextRun.Content {
						size := utf8.RuneLen(r)
						if size > 1 {
							rUtf16 := utf16.Encode([]rune{r})
							rUtf16Size := len(rUtf16)
							size = rUtf16Size
						}
						startIndex := index
						index += int64(size)
						chars = append(chars, char{startIndex, index, r})
					}
				}
			}
		}
	}
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
	// read client secret
	bytes, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Failed to read credentials: %v", err)
	}

	// initialize config for client authorization
	config, err := google.ConfigFromJSON(bytes, scope)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// authorize client (OAuth2)
	client := authorizeClient(config)

	// create docs service -- later use api
	docsService, err := docs.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Docs service: %v", err)
	}

	// do stuff!
	subscribe(docsService)
}
