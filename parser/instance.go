package parser

import (
	"GDocs-Syntax-Highlighter/request"
	"GDocs-Syntax-Highlighter/style"
	"log"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/docs/v1"
)

const (
	codeInstanceStart = "<code>"  // required tag to denote start of code instance
	codeInstanceEnd   = "</code>" // required tag to denote end of code instance
	configStart       = "<conf>"  // required tag to denote start of config
	configEnd         = "</conf>" // required tag to denote end of config

	// Optional directive to specify if the code should be formatted.
	// Note that formatting is not highlighting.
	// If not present, the code will never be formatted.
	// If present, the code is formatted every time the user bolds this config directive.
	formatDirective = "#format"
)

var (
	// Optional directive to specify the language of the code.
	// If not set, #lang=<go> is assumed by default.
	configLangRegex = regexp.MustCompile("^#lang=(\\w+)$")

	// Optional directive to specify the font of the code.
	// If not set, #font=<Courier New> is assumed by default.
	fontRegex = regexp.MustCompile("^#font=([\\w_]+)$")

	// Optional directive to specify the font size of the code.
	// If not set, #size=<11> is assumed by default.
	fontSizeRegex = regexp.MustCompile("^#size=(\\d+(\\.\\d+)?)$")
)

// CodeInstance describes a section in the Google Doc
// that has a config and code fragment.
type CodeInstance struct {
	builder          strings.Builder // string builder for code body
	foundConfigStart bool            // whether the config start tag was found
	foundConfigEnd   bool            // whether the config end tag was found
	Code             string          // the code as text
	Font             *string         // font
	FontSize         *float64        // font size
	Lang             *style.Language // the coding language
	StartIndex       int64           // start index of code
	EndIndex         int64           // end index of code
	Format           *style.Format   // whether we are being requested to format the code
}

// Range is a helper function to get the *docs.Range
// for a particular code instance.
func (c *CodeInstance) Range() *docs.Range {
	return request.GetRange(c.StartIndex, c.EndIndex)
}

// Sets default values if unset.
func (c *CodeInstance) setDefaults() {
	if c.Lang == nil {
		defaultLang := style.GetDefaultLanguage()
		c.Lang = &defaultLang
	}
	if c.Format == nil {
		c.Format = &style.Format{}
	}
	if c.Font == nil {
		defaultFont := style.DefaultFont
		c.Font = &defaultFont
	}
	if c.FontSize == nil {
		defaultSize := style.DefaultFontSize
		c.FontSize = &defaultSize
	}
}

// Checks for header tags/directives in a particular
// string that is located in a *docs.ParagraphElement.
func (c *CodeInstance) checkForHeader(s string, par *docs.ParagraphElement) {
	// search for start of config tags
	if !c.foundConfigStart {
		if strings.EqualFold(s, configStart) {
			c.foundConfigStart = true
		}
		return
	}

	// check for end of config
	if strings.EqualFold(s, configEnd) {
		c.foundConfigEnd = true
		c.StartIndex = par.EndIndex
		c.setDefaults()
		return
	}

	// check for format directive (and bolded)
	if c.Format == nil && strings.EqualFold(s, formatDirective) {
		formatStart, formatEnd := getUTF16SubstrIndices(formatDirective, par.TextRun.Content, par.StartIndex)
		c.Format = &style.Format{
			Bold:       par.TextRun.TextStyle.Bold,
			StartIndex: formatStart,
			EndIndex:   formatEnd,
		}
		return
	}

	// check for language directive
	if c.Lang == nil {
		if res := configLangRegex.FindStringSubmatch(s); len(res) == 2 {
			if lang, ok := style.GetLanguage(res[1]); ok {
				c.Lang = &lang
			} else {
				// TODO: maybe add a comment to the Google Doc
				// in the future to notify of an invalid language name
				log.Printf("Unknown language: `%s`\n", res[1])
			}
			return
		}
	}

	// check for font
	if c.Font == nil {
		if res := fontRegex.FindStringSubmatch(s); len(res) == 2 {
			if font, ok := style.GetFont(res[1]); ok {
				c.Font = &font
			} else {
				// TODO: maybe add a comment to the Google Doc
				// in the future to notify of an invalid language name
				log.Printf("Unknown font: `%s`\n", res[1])
			}
			return
		}
	}

	// check for font size
	if c.FontSize == nil {
		if res := fontSizeRegex.FindStringSubmatch(s); len(res) == 3 {
			float, err := strconv.ParseFloat(res[1], 64)
			if err != nil {
				log.Printf("Failed to parse font size `%s` into float64: %s\n", res[1], err)
			} else {
				c.FontSize = &float // if it is 0, will default to 1
			}
			return
		}
	}

	log.Printf("Unexpected config token: `%s`\n", s)
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
					italics := par.TextRun.TextStyle.Italic

					if cur == nil || !cur.foundConfigEnd {
						// iterate over each word
						for _, s := range strings.Fields(content) {
							// note: all tags must be in italics to separate them
							// from any collision with the code body
							if !italics {
								continue // ignore non-italics
							}

							// have not found start of instance yet so check for start symbol
							if cur == nil {
								if strings.EqualFold(s, codeInstanceStart) {
									cur = &CodeInstance{}
								}
								continue
							}

							cur.checkForHeader(s, par)
						}
						continue
					}

					// check for footer/end symbol
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
