package runner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// A request to the Go Playground.
type goPlaygroundRequest struct {
	Body string `json:"body"`
}

// goPlaygroundResponse is a response from the Go Playground.
type goPlaygroundResponse struct {
	Errors             string              `json:"errors"`
	GoPlaygroundEvents []goPlaygroundEvent `json:"events"`
	Status             int                 `json:"status"`
	IsTest             bool                `json:"istest"`
	TestsFailed        int                 `json:"testsfailed"`
}

// goPlaygroundEvent is an output event.
type goPlaygroundEvent struct {
	Message string `json:"message"`
	Kind    string `json:"kind"`
	Delay   int    `json:"delay"`
}

// RunResult represents the result of running a program.
type RunResult struct {
	Output string // for now combing stdout, stderr
	Errors string
	Status int
}

// RunGo runs Go using Go Playground's server.
func RunGo(program string) (*RunResult, error) {
	// marshal payload
	payload, err := json.Marshal(goPlaygroundRequest{program})
	if err != nil {
		return nil, err
	}

	// send request to Go Playground
	resp, err := http.Post("https://play.golang.org/compile", "text/plain", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal response
	var goResp goPlaygroundResponse
	err = json.Unmarshal(body, &goResp)
	if err != nil {
		return nil, nil
	}

	// combine stderr and stdout for now
	var b strings.Builder
	for _, event := range goResp.GoPlaygroundEvents {
		_, err = b.WriteString(event.Message)
		if err != nil {
			return nil, err
		}
	}

	return &RunResult{
		Output: b.String(),
		Errors: goResp.Errors,
		Status: goResp.Status,
	}, nil
}
