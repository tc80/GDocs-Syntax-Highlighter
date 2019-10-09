package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	// encode the Token into json
	json.NewEncoder(file).Encode(token)
}

// for testing now
func subscribe(docsService *docs.Service) {
	docID := ""
	doc, err := docsService.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Failed to get doc: %v", err)
	}
	fmt.Printf("Document: %v\n", doc.Title)
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

	// create docs service
	docsService, err := docs.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Docs service: %v", err)
	}

	// do stuff!
	subscribe(docsService)
}
