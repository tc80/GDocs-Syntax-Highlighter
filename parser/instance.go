package parser

import (
	"GDocs-Syntax-Highlighter/request"
	"GDocs-Syntax-Highlighter/style"
	"log"
	"strconv"
	"strings"

	"google.golang.org/api/docs/v1"
)

// ConfigSegment represents a header/footer, which
// if where config directives live.
type ConfigSegment struct {
	StartIndex int64
	EndIndex   int64
}

// CodeInstance describes a section in the Google Doc
// that has a config and code fragment.
type CodeInstance struct {
	toUTF16    map[int]int64             // maps the indices of the zero-based utf8 rune in Code to utf16 rune indices+start utf16 offset
	Segments   map[string]*ConfigSegment // headers and footer IDs -> config segment
	Code       string                    // the code as text
	Theme      *string                   // theme
	Font       *string                   // font
	FontSize   *float64                  // font size
	Lang       *style.Language           // the coding language
	StartIndex *int64                    // utf16 start index of code
	EndIndex   *int64                    // utf16 end index of code
	Shortcuts  *bool                     // whether shortcuts are enabled
	Format     *style.Format             // whether we are being requested to format the code
}

// GetRange gets the *docs.Range
// for a particular code instance.
func (c *CodeInstance) GetRange() *docs.Range {
	return request.GetRange(*c.StartIndex, *c.EndIndex, "")
}

// GetTheme gets the *style.Theme for a particular code instance.
// Note that the language and theme fields must be valid.
func (c *CodeInstance) GetTheme() *style.Theme {
	return c.Lang.Themes[*c.Theme]
}

// UpdateCode gets the []*docs.Request to delete the existing
// code range and replace it with a new string Code.
// It does not update the indices.
func (c *CodeInstance) UpdateCode() []*docs.Request {
	return []*docs.Request{
		// need to ignore the newline character at the end of the segment so we use EndIndex-1
		request.Delete(request.GetRange(*c.StartIndex, *c.EndIndex-1, "")),
		request.Insert(c.Code, *c.StartIndex),
	}
}

// Sets default values if unset.
// Does not set start/end indices.
func (c *CodeInstance) setDefaults() {
	if c.Lang == nil {
		c.Lang = style.GetDefaultLanguage()
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
	if c.Theme == nil {
		defaultTheme := style.DefaultTheme
		c.Theme = &defaultTheme
	}
	if c.Shortcuts == nil {
		defaultShortcuts := style.DefaultShortcutSetting
		c.Shortcuts = &defaultShortcuts
	}
	if c.toUTF16 == nil {
		c.toUTF16 = make(map[int]int64)
	}
}

// Checks for config directives in a particular
// string that is located in a *docs.ParagraphElement.
func (c *CodeInstance) checkForConfig(s, segmentID string, par *docs.ParagraphElement) {
	// check for format (must be bolded)
	if c.Format == nil && strings.EqualFold(s, style.FormatDirective) {
		formatStart, formatEnd := getUTF16SubstrIndices(style.FormatDirective, par.TextRun.Content, par.StartIndex)
		c.Format = &style.Format{
			Bold:       par.TextRun.TextStyle.Bold,
			StartIndex: formatStart,
			EndIndex:   formatEnd,
			SegmentID:  segmentID,
		}
		return
	}

	// check for shortcuts
	if c.Shortcuts == nil {
		if res := style.ShortcutsRegex.FindStringSubmatch(s); len(res) == 2 {
			enabled := res[1] == "enabled"
			c.Shortcuts = &enabled
			return
		}
	}

	// check for language
	if c.Lang == nil {
		if res := style.LangRegex.FindStringSubmatch(s); len(res) == 2 {
			if l, ok := style.GetLanguage(res[1]); ok {
				c.Lang = l
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
		if res := style.FontRegex.FindStringSubmatch(s); len(res) == 2 {
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
		if res := style.FontSizeRegex.FindStringSubmatch(s); len(res) == 3 {
			float, err := strconv.ParseFloat(res[1], 64)
			if err != nil {
				log.Printf("Failed to parse font size `%s` into float64: %s\n", res[1], err)
			} else {
				c.FontSize = &float // if it is 0, will default to 1
			}
			return
		}
	}

	// check for theme
	if c.Theme == nil {
		if res := style.ThemeRegex.FindStringSubmatch(s); len(res) == 2 {
			if theme, ok := style.GetTheme(res[1]); ok {
				c.Theme = &theme
			} else {
				// TODO: maybe add a comment to the Google Doc
				// in the future to notify of an invalid language name
				log.Printf("Unknown theme: `%s`\n", res[1])
			}
			return
		}
	}

	log.Printf("Unexpected config token: `%s`\n", s)
}

// GetCodeInstance gets the config and instances of code and that
// will be processed in a Google Doc.
// Note that directives split into multiple text runs are ignored.
func GetCodeInstance(doc *docs.Document) *CodeInstance {
	c := new(CodeInstance)
	c.Segments = make(map[string]*ConfigSegment)

	// check for config in Google Doc headers
	for _, h := range doc.Headers {
		for _, elem := range h.Content {
			if elem.Paragraph != nil {
				for _, par := range elem.Paragraph.Elements {
					if par.TextRun != nil {
						if seg, ok := c.Segments[h.HeaderId]; ok {
							seg.EndIndex = par.EndIndex
						} else {
							c.Segments[h.HeaderId] = &ConfigSegment{par.StartIndex, par.EndIndex}
						}
						for _, s := range strings.Fields(par.TextRun.Content) {
							c.checkForConfig(s, h.HeaderId, par)
						}
					}
				}
			}
		}
	}

	// check for config in Google Doc footers
	for _, f := range doc.Footers {
		for _, elem := range f.Content {
			if elem.Paragraph != nil {
				for _, par := range elem.Paragraph.Elements {
					if par.TextRun != nil {
						if seg, ok := c.Segments[f.FooterId]; ok {
							seg.EndIndex = par.EndIndex
						} else {
							c.Segments[f.FooterId] = &ConfigSegment{par.StartIndex, par.EndIndex}
						}
						for _, s := range strings.Fields(par.TextRun.Content) {
							c.checkForConfig(s, f.FooterId, par)
						}
					}
				}
			}
		}
	}

	// concatenate Google Doc body
	var b strings.Builder
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					if c.StartIndex == nil {
						c.StartIndex = &par.StartIndex
					}
					c.EndIndex = &par.EndIndex
					_, err := b.WriteString(par.TextRun.Content)
					check(err)
				}
			}
		}
	}
	c.Code = b.String()

	// set defaults
	c.setDefaults()

	return c
}
