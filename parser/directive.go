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

var (
	// formatDirective is an optional directive to specify if the code should be formatted.
	// Note that formatting is not highlighting.
	// If not present, the code will never be formatted.
	// If present, the code is formatted every time the user underlines this config directive.
	formatDirective = "#format"

	// runDirective is an optional directive to specify if the code should be run in a sandbox.
	// If not present, the code will never be run.
	// If present, the code is run every time the user underlines this config directive.
	runDirective = "#run"

	// FontRegex is an optional directive to specify the font of the code.
	// If not set, #font=courier_new is assumed by default.
	fontDirectiveRegex = regexp.MustCompile("^#font=([\\w_]+)$")

	// FontSizeRegex is an optional directive to specify the font size of the code.
	// If not set, #size=11 is assumed by default.
	fontSizeDirectiveRegex = regexp.MustCompile("^#size=(\\d+(\\.\\d+)?)$")

	// LangRegex is the regex for the optional directive
	// to specify the language of the code.
	// If not set, #lang=go is assumed by default.
	langDirectiveRegex = regexp.MustCompile("^#lang=([\\w_]+)$")

	// ShortcutsRegex is an optional directive to specify if shortcuts are enabled.
	// By default, shortcuts are disabled.
	shortcutsDirectiveRegex = regexp.MustCompile("^#shortcuts=(enabled|disabled)$")

	// ThemeRegex is an optional directive to specify the theme of the code.
	// If not set, #theme=dark is assumed by default.
	themeDirectiveRegex = regexp.MustCompile("^#theme=([\\w_]+)$")
)

// UnderlinedDirective describes a directive that is underlined
// as well as the UTF16 indices of the directive (to un-underline itself).
type UnderlinedDirective struct {
	Underlined bool   // if underlined, do something and then un-underline the directive
	SegmentID  string // segment ID
	StartIndex int64  // start index of directive
	EndIndex   int64  // end index of directive
}

// GetRange gets the *docs.Range
// for a particular UnderlinedDirective.
func (u *UnderlinedDirective) GetRange() *docs.Range {
	return request.GetRange(u.StartIndex, u.EndIndex, u.SegmentID)
}

// Checks for config directives in a particular
// string that is located in a *docs.ParagraphElement.
func (c *CodeInstance) checkForDirectives(s, segmentID string, par *docs.ParagraphElement) {
	// check for format (must be underlined)
	if c.Format == nil && strings.EqualFold(s, formatDirective) {
		formatStart, formatEnd := getUTF16SubstrIndices(formatDirective, par.TextRun.Content, par.StartIndex)
		c.Format = &UnderlinedDirective{
			Underlined: par.TextRun.TextStyle.Underline,
			StartIndex: formatStart,
			EndIndex:   formatEnd,
			SegmentID:  segmentID,
		}
		return
	}

	// check for run (must be underlined)
	if c.Run == nil && strings.EqualFold(s, runDirective) {
		runStart, runEnd := getUTF16SubstrIndices(runDirective, par.TextRun.Content, par.StartIndex)
		c.Run = &UnderlinedDirective{
			Underlined: par.TextRun.TextStyle.Underline,
			StartIndex: runStart,
			EndIndex:   runEnd,
			SegmentID:  segmentID,
		}
		return
	}

	// check for shortcuts
	if c.Shortcuts == nil {
		if res := shortcutsDirectiveRegex.FindStringSubmatch(s); len(res) == 2 {
			enabled := res[1] == "enabled"
			c.Shortcuts = &enabled
			return
		}
	}

	// check for language
	if c.Lang == nil {
		if res := langDirectiveRegex.FindStringSubmatch(s); len(res) == 2 {
			if l, ok := style.GetLanguage(res[1]); ok {
				c.Lang = l
			} else {
				// TODO: highlight invalid directive
				log.Printf("Unknown language: `%s`\n", res[1])
			}
			return
		}
	}

	// check for font
	if c.Font == nil {
		if res := fontDirectiveRegex.FindStringSubmatch(s); len(res) == 2 {
			if font, ok := style.GetFont(res[1]); ok {
				c.Font = &font
			} else {
				// TODO: highlight invalid directive
				log.Printf("Unknown font: `%s`\n", res[1])
			}
			return
		}
	}

	// check for font size
	if c.FontSize == nil {
		if res := fontSizeDirectiveRegex.FindStringSubmatch(s); len(res) == 3 {
			float, err := strconv.ParseFloat(res[1], 64)
			if err != nil {
				// TODO: highlight invalid directive
				log.Printf("Failed to parse font size `%s` into float64: %s\n", res[1], err)
			} else {
				c.FontSize = &float // if it is 0, will default to 1
			}
			return
		}
	}

	// check for theme
	if c.Theme == nil {
		if res := themeDirectiveRegex.FindStringSubmatch(s); len(res) == 2 {
			if theme, ok := style.GetTheme(res[1]); ok {
				c.Theme = &theme
			} else {
				// TODO: highlight invalid directive
				log.Printf("Unknown theme: `%s`\n", res[1])
			}
			return
		}
	}

	// TODO: highlight invalid directive
	log.Printf("Unexpected config token: `%s`\n", s)
}
