package parser

import (
	"GDocs-Syntax-Highlighter/request"
	"GDocs-Syntax-Highlighter/style"
	"regexp"
	"strings"
	"unicode/utf8"

	"google.golang.org/api/docs/v1"
)

// Instance of parserInput for parsing a range.
type rangeInput struct {
	pos   int
	runes string
}

// Output from a rangeInput parser.
type rangeOutput struct {
	result    string
	rangeType *style.Range
}

// Gets the current rune and its size.
func (in rangeInput) current() (*rune, int) {
	if in.pos >= len(in.runes) {
		return nil, 0
	}
	r, size := utf8.DecodeRuneInString(in.runes[in.pos:])
	if r == utf8.RuneError {
		panic("invalid rune")
	}
	return &r, size
}

// Advances to the next rune based on the previous rune's size.
func (in rangeInput) advance(size int) parserInput {
	return rangeInput{in.pos + size, in.runes}
}

// Remove a string of characters at a utf8 index.
type removeRange struct {
	index    int
	utf8Size int
}

// RemoveRanges removes the ranges from the instance's Code
// string property and returns the list of requests to highlight them.
func (c *CodeInstance) RemoveRanges(t *style.Theme) (reqs []*docs.Request) {
	// create parser for each range
	var rangeParsers []parser
	for _, r := range t.Ranges {
		rangeParsers = append(rangeParsers, expectRange(r))
	}

	var removeRanges []removeRange      // ranges to be removed
	utf16Offsets := make(map[int]int64) // add utf16 offset at certain utf8 indices in the sanitized string
	var utf8Offset int                  // utf8 offset in sanitized string
	var utf8MultiRangeStartIndex *int   // if multiple ranges in a row, utf8 start index in sanitized string of first range
	in := rangeInput{runes: c.Code}
	for r, size := in.current(); r != nil; r, size = in.current() {
		out := selectAny(rangeParsers)(in)
		if out.result != nil {
			// if range found, consume it and remove from string
			rOutput := out.result.(rangeOutput)
			utf8StartIndex := in.pos
			utf8Size := len(rOutput.result)
			utf16Size := GetUtf16StringSize(rOutput.result)
			removeRanges = append(removeRanges, removeRange{utf8StartIndex, utf8Size})

			// update offset map to recreate utf8 -> utf16 map in sanitized string
			if utf8MultiRangeStartIndex != nil {
				// there is at least one range directly before this one
				// so update the first range's utf8 -> utf16 index mapping in the sanitized string
				utf16Offsets[*utf8MultiRangeStartIndex] = utf16Offsets[*utf8MultiRangeStartIndex] + utf16Size
			} else {
				// no range directly before this, so create new utf8 -> utf16 index mapping
				utf8RangeStartIndex := utf8StartIndex + utf8Offset
				utf16Offsets[utf8RangeStartIndex] = utf16Size
				utf8MultiRangeStartIndex = &utf8RangeStartIndex
			}
			utf8Offset -= utf8Size

			// create request to update range's color using utf16 indices
			utf16OffsetStartIndex := c.toUTF16[utf8StartIndex]
			utf16Range := request.GetRange(utf16OffsetStartIndex, utf16OffsetStartIndex+utf16Size, "")
			reqs = append(reqs, request.UpdateForegroundColor(rOutput.rangeType.Color, utf16Range))

			in = out.remaining.(rangeInput)
			continue
		}
		// failed to parse range, advance
		in = in.advance(size).(rangeInput)
		utf8MultiRangeStartIndex = nil
	}

	if len(removeRanges) == 0 {
		return
	}

	var sanitized strings.Builder
	var cur removeRange

	// remove ranges from Code
	for len(removeRanges) > 0 {
		start := cur.index + cur.utf8Size
		cur, removeRanges = removeRanges[0], removeRanges[1:] // pop from slice
		_, err := sanitized.WriteString(c.Code[start:cur.index])
		check(err)
	}
	_, err := sanitized.WriteString(c.Code[cur.index+cur.utf8Size:])
	check(err)

	// update Code (removed ranges)
	c.Code = sanitized.String()

	// update zero-based utf8 -> offset utf16 index mapping
	utf16Index := *c.StartIndex
	c.toUTF16 = make(map[int]int64)
	var r rune
	for i, utf8Width := 0, 0; i < len(c.Code); i += utf8Width {
		if o, ok := utf16Offsets[i]; ok {
			utf16Index += o // add utf16 offset for range removal
		}
		r, utf8Width = utf8.DecodeRuneInString(c.Code[i:])
		c.toUTF16[i] = utf16Index
		utf16Index += GetUtf16RuneSize(r)
	}
	return
}

// Replace gets the requests to replace all matches of a regex with a particular string.
// It also updates the instance.Code and EndIndex respectively.
func (c *CodeInstance) Replace(s *style.Shortcut) (reqs []*docs.Request) {
	for {
		if res := s.Regex.FindStringSubmatchIndex(c.Code); res != nil {
			utf8DelStart, utf8DelEnd := res[0], res[1]
			utf16DelStart, utf16DelEnd := getUTF16SubstrIndices(c.Code[utf8DelStart:utf8DelEnd], c.Code, *c.StartIndex)

			// delete target and insert replacement string
			utf16DelRange := request.GetRange(utf16DelStart, utf16DelEnd, "")
			reqs = append(reqs, request.Delete(utf16DelRange))
			reqs = append(reqs, request.Insert(s.Replace, utf16DelStart))

			// update end index for utf16 difference
			utf16InsSize := GetUtf16StringSize(s.Replace)
			newEndIndex := *c.EndIndex + utf16InsSize - (utf16DelEnd - utf16DelStart)
			c.EndIndex = &newEndIndex

			// replace c.Code
			c.Code = c.Code[:utf8DelStart] + s.Replace + c.Code[utf8DelEnd:]
			continue
		}
		return
	}
}

// Highlight gets the requests to highlight all matches of a regex with a particular color.
func (c *CodeInstance) Highlight(r *regexp.Regexp, color *docs.Color, segmentID string) (reqs []*docs.Request) {
	if results := r.FindAllStringSubmatchIndex(c.Code, -1); results != nil {
		for _, res := range results {
			utf8Start, utf8End := res[0], res[1]
			utf16Size := GetUtf16StringSize(c.Code[utf8Start:utf8End])
			utf16StartOffset := c.toUTF16[utf8Start]
			utf16Range := request.GetRange(utf16StartOffset, utf16StartOffset+utf16Size, segmentID)
			reqs = append(reqs, request.UpdateForegroundColor(color, utf16Range))
		}
	}
	return
}
