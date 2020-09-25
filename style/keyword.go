package style

import (
	"regexp"

	"google.golang.org/api/docs/v1"
)

// Keyword represents a language keyword
// and the color it is highlighted with (for a theme).
type Keyword struct {
	Regex *regexp.Regexp
	Color *docs.Color
}

var (
	// Note that some of the following Go regexes are taken/inspired from the VSCode language files found here:
	// https://github.com/microsoft/vscode/blob/master/extensions/go/syntaxes/go.tmLanguage.json
	//
	// TODO: update regexes to use capturing groups to replace positive lookaheads (since Go doesn't support it)
	// store some metadata about which regex group to use to highlight?
	go1 = regexp.MustCompile("\\b(break|case|continue|default|defer|else|fallthrough|for|go|goto|if|range|return|select|switch)\\b")
	go2 = regexp.MustCompile("\\b(chan|const|func|interface|map|struct|true|false|nil|iota|package|type|import|var)\\b")
	go3 = regexp.MustCompile("\\b(bool|byte|error|(complex(64|128)|float(32|64)|u?int(8|16|32|64)?)|rune|string|uintptr)\\b")
	go4 = regexp.MustCompile("\\b(append|cap|close|complex|copy|delete|imag|len|make|new|panic|print|println|real|recover)\\b")
	go5 = regexp.MustCompile("\\b\\d+\\b")
)
