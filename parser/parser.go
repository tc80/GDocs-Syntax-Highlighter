package parser

import (
	"GDocs-Syntax-Highlighter/request"
	"GDocs-Syntax-Highlighter/style"
	"fmt"
	"strings"
	"unicode/utf8"

	"google.golang.org/api/docs/v1"
)

// Represents a parser.
type parser func(parserInput) parserOutput

// Functions to get the current rune the parser is processing
// and to advance the rune stream.
type parserInput interface {
	current() (*rune, int)   // return current rune, its size
	advance(int) parserInput // advance based on rune size
}

// The parsed result and the remaining stream.
type parserOutput struct {
	result    interface{}
	remaining parserInput
}

type rangeResult struct {
	start int
	end   int
	color *docs.Color
}

// Provides input for parsing a range.
type rangeInput struct {
	pos   int
	runes string
}

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

// Denotes a parse success.
func success(result interface{}, input parserInput) parserOutput {
	return parserOutput{result, input}
}

// Denotes a parse failure.
func fail() parserOutput {
	return parserOutput{nil, nil}
}

// Enforces a property is non-nil.
func check(e interface{}) {
	if e != nil {
		panic(fmt.Sprintf("check fail: %v", e))
	}
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

	var removeRanges []removeRange
	utf16Offsets := make(map[int]int64) // add utf16 offset at certain utf8 indices in the sanitized string
	var utf8Offset int
	in := rangeInput{runes: c.Code}
	for r, size := in.current(); r != nil; r, size = in.current() {
		out := selectAny(rangeParsers)(in)
		if out.result != nil {
			// if range found, consume it and remove from string
			rOutput := out.result.(rangeOutput)
			utf8StartIndex := in.pos
			utf8Size := len(rOutput.result)
			utf16Size := GetUtf16StringSize(rOutput.result)
			utf16Offsets[utf8StartIndex+utf8Offset] = utf16Size
			utf8Offset -= utf8Size
			removeRanges = append(removeRanges, removeRange{utf8StartIndex, utf8Size})

			// create request to update range's color using utf16 indices
			utf16OffsetStartIndex := c.toUTF16[utf8StartIndex]
			utf16Range := request.GetRange(utf16OffsetStartIndex, utf16OffsetStartIndex+utf16Size)
			reqs = append(reqs, request.UpdateForegroundColor(rOutput.rangeType.Color, utf16Range))

			in = out.remaining.(rangeInput)
			continue
		}
		in = in.advance(size).(rangeInput)
	}

	if len(removeRanges) == 0 {
		return
	}

	var sanitized strings.Builder
	var cur removeRange

	// if something to remove, remove range from
	for len(removeRanges) > 0 {
		start := cur.index + cur.utf8Size
		cur, removeRanges = removeRanges[0], removeRanges[1:] // pop from slice
		sanitizedStr := c.Code[start:cur.index]
		_, err := sanitized.WriteString(sanitizedStr)
		check(err)
	}
	_, err := sanitized.WriteString(c.Code[cur.index+cur.utf8Size:])
	check(err)

	// update code (removed ranges)
	c.Code = sanitized.String()

	// update zero-based utf8 -> offset utf16 index mapping
	utf16Index := c.StartIndex
	c.toUTF16 = make(map[int]int64)
	var r rune
	for i, utf8Width := 0, 0; i < len(c.Code); {
		if o, ok := utf16Offsets[i]; ok {
			utf16Index += o // add utf16 offset for range removal
		}

		r, utf8Width = utf8.DecodeRuneInString(c.Code[i:])
		c.toUTF16[i] = utf16Index
		fmt.Println(i, " -> ", utf16Index)

		utf16Index += GetUtf16RuneSize(r)
		i += utf8Width
	}

	fmt.Println(c.EndIndex)
	return
}

// Selects the first parser in a slice of
// parsers that successfully parses the input.
func selectAny(parsers []parser) parser {
	return func(in parserInput) parserOutput {
		for _, p := range parsers {
			if out := p(in); out.result != nil {
				return out
			}
		}
		return fail() // all parsers failed
	}
}

// Parser for a symbol range.
// The parser returns
func expectRange(r *style.Range) parser {
	return func(in parserInput) parserOutput {
		// check for start symbol
		out := expectString(r.StartSymbol)(in)
		if out.result == nil {
			return fail()
		}
		in = out.remaining
		var b strings.Builder
		_, err := b.WriteString(r.StartSymbol)
		check(err)

		// search until end symbol or end
		out = searchUntil(expectString(r.EndSymbol))(in)
		s := out.result.(search)
		_, err = b.WriteString(s.consumed)
		check(err)

		// if end symbol found, add to builder
		if s.result != nil {
			_, err = b.WriteString(r.EndSymbol)
			check(err)
		}
		in = out.remaining
		return success(rangeOutput{b.String(), r}, in)
	}
}

// Represents a search
type search struct {
	consumed string      // string of consumed runes while searching
	result   interface{} // if the parser parsed something, the result would be here
}

// Parser that keep consuming all runes until the parser is successful
// or the end is reached. It returns a search struct.
func searchUntil(p parser) parser {
	return func(in parserInput) parserOutput {
		var consumed strings.Builder
		out := p(in)
		for ; out.result == nil; out = p(in) {
			out = expectRune(anyRune())(in)
			if out.result == nil {
				// reached end, parser did not find anything
				return success(search{consumed.String(), nil}, in)
			}
			_, err := consumed.WriteRune(out.result.(rune))
			check(err)
			in = out.remaining
		}
		// parser consumed something, so return
		return success(search{consumed.String(), out.result}, out.remaining)
	}
}

// Expects an exact string, rune-by-rune.
// If success, parser returns the string.
func expectString(s string) parser {
	return func(in parserInput) parserOutput {
		for _, r := range s {
			out := expectRune(isRune(r))(in)
			if out.result == nil {
				return fail()
			}
			in = out.remaining
		}
		return success(s, in)
	}
}

// Expects a given rune based on a boolean function.
// If success, the parser returns the *rune.
func expectRune(ok isRuneFunc) parser {
	return func(in parserInput) parserOutput {
		r, size := in.current()
		if r == nil || !ok(*r) {
			return fail()
		}
		return success(*r, in.advance(size))
	}
}
