package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
)

const (
	scope           = docs.DriveScope         // needed for editing GDrive files
	stateToken      = "state-token"           // used for requesting a new token
	credentialsPath = "auth/credentials.json" // client secret
	tokenPath       = "auth/token.json"       // token path, needs to change if scope changes
)

// Authorizes the client with an API token.
func authorizeClient(config *oauth2.Config) (*http.Client, error) {
	token, err := checkForToken()
	if err != nil {
		log.Println("Unable to locate local token, attempting to get token from web.")
		token, err = requestNewToken(config)
		if err != nil {
			return nil, err
		}
		cacheToken(token)
	}
	return config.Client(context.Background(), token), nil
}

// Request a new token from the Docs API.
func requestNewToken(config *oauth2.Config) (*oauth2.Token, error) {
	// get authorization code
	log.Printf("Enter auth code from: \n%v\n", config.AuthCodeURL(stateToken, oauth2.AccessTypeOffline))
	var auth string
	_, err := fmt.Scan(&auth)
	if err != nil {
		return nil, errors.New("Failed to scan auth code: " + err.Error())
	}

	// get new token using auth code, passing empty context (same as TODO())
	token, err := config.Exchange(oauth2.NoContext, auth)
	if err != nil {
		return nil, errors.New("Failed to get token: " + err.Error())
	}
	return token, nil
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
		log.Printf("\nFailed to write token to file: %v", err)
		return
	}

	// encode the token into json
	json.NewEncoder(file).Encode(token)
}

// GetAuthorizedClient gets the authorized (OAuth2) http Client
func GetAuthorizedClient() (*http.Client, error) {
	// read client secret
	bytes, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, errors.New("Failed to read credentials: " + err.Error())
	}

	// initialize config for client authorization
	config, err := google.ConfigFromJSON(bytes, scope)
	if err != nil {
		return nil, errors.New("Failed to parse config: " + err.Error())
	}

	// authorize client (OAuth2)
	return authorizeClient(config)
}
