package parser

import (
	"GDocs-Syntax-Highlighter/style"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/api/docs/v1"
)

const (
	codeInstanceStart = "<code>"    // required tag to denote start of code instance
	codeInstanceEnd   = "</code>"   // required tag to denote end of code instance
	configStart       = "<config>"  // required tag to denote start of config
	configEnd         = "</config>" // required tag to denote end of config

	// Optional directive to specify if the code should be formatted.
	// Note that formatting is not highlighting.
	// If not present, the code will never be formatted.
	// If present, the code is formatted every time the user bolds this config directive.
	formatDirective = "#format"
)

var (
	// Optional directive to specify the language of the code.
	// If not set, #lang=go is assumed by default.
	configLangRegex = regexp.MustCompile("^#lang=(\\w+)$")
)

// CodeInstance describes a section in the Google Doc
// that has a config and code fragment.
type CodeInstance struct {
	builder          strings.Builder // string builder for code body
	foundConfigStart bool            // whether the config start tag was found
	foundConfigEnd   bool            // whether the config end tag was found
	Code             string          // the code as text
	Lang             *style.Language // the coding language
	StartIndex       int64           // start index of code
	EndIndex         int64           // end index of code
	Format           *Format         // whether we are being requested to format the code
}

// Format describes whether we will format the code (if the directive is bolded)
// as well as the UTF16 indices of the directive (to unbold itself).
type Format struct {
	Bold       bool  // if bolded, format the code and then unbold the directive
	StartIndex int64 // start index of directive
	EndIndex   int64 // end index of directive

}

// GetCodeInstances gets the instances of code that will be processed in
// a Google Doc. Each instance will be surrounded with <code> and </code> tags, as
// well as a header containing info for configuration with <config> and </config> tags.
func GetCodeInstances(doc *docs.Document) []*CodeInstance {
	var instances []*CodeInstance
	var cur *CodeInstance

	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					content := par.TextRun.Content
					italics, bold := par.TextRun.TextStyle.Italic, par.TextRun.TextStyle.Bold

					if cur == nil || !cur.foundConfigEnd {
						// iterate over each word
						for _, str := range strings.Fields(content) {
							if !italics {
								continue // ignore non-italics
							}

							// have not found start of instance yet so check for start symbol
							// note: all tags must be in italics to separate them from any collision with the code body
							if cur == nil {
								if strings.EqualFold(str, codeInstanceStart) {
									cur = &CodeInstance{}
								}
								continue
							}

							// search for start of config tags
							if !cur.foundConfigStart {
								if strings.EqualFold(str, configStart) {
									cur.foundConfigStart = true
								}
								continue
							}

							// search for config tags/directives

							// check for end of config
							if strings.EqualFold(str, configEnd) {
								cur.foundConfigEnd = true
								cur.StartIndex = par.EndIndex
								// set defaults if necessary
								if cur.Lang == nil {
									defaultLang := style.GetDefaultLanguage()
									cur.Lang = &defaultLang
								}
								if cur.Format == nil {
									cur.Format = &Format{}
								}
								continue
							}

							// check for format directive (and bolded)
							if cur.Format == nil && strings.EqualFold(str, formatDirective) {
								formatStart, formatEnd := getUTF16SubstrIndices(formatDirective, content, par.StartIndex)
								cur.Format = &Format{
									Bold:       bold,
									StartIndex: formatStart,
									EndIndex:   formatEnd,
								}
								continue
							}

							// check for language directive
							if cur.Lang == nil {
								if res := configLangRegex.FindStringSubmatch(str); len(res) == 2 {
									if lang, ok := style.GetLanguage(res[1]); ok {
										cur.Lang = &lang
									} else {
										// TODO: maybe add a comment to the Google Doc
										// in the future to notify of an invalid language name
										fmt.Printf("Unknown language: `%s`\n", res[1])
									}
									continue
								}
							}
							fmt.Printf("Unexpected config token: `%s`\n", str)
							continue
						}
						continue
					}

					// check for end symbol
					if italics && strings.EqualFold(strings.TrimSpace(content), codeInstanceEnd) {
						cur.Code = cur.builder.String()
						instances = append(instances, cur)
						cur = nil
						continue
					}

					// write untrimmed body content, update end index
					cur.builder.WriteString(content)
					cur.EndIndex = par.EndIndex
				}
			}
		}
	}
	return instances
}
