package main

import (
	"bytes"
	"fmt"
	"unicode/utf16"
)

// A rune and its respective utf16 start and end indices
type char struct {
	index   int64 // the utf16 inclusive start index of the rune
	size    int64 // the size of the rune in utf16 units
	content rune  // the rune
}

// A word as a string and its respective utf16 start and end indices
type word struct {
	index   int64  // the utf16 inclusive start index of the string
	size    int64  // the size of the word in utf16 units
	content string // the string
}

// Comment
type comment struct {
	startSymbol string
	endSymbol   string
}

type parserInput interface {
	current() *char
	advance() parserInput
}

type parserOutput struct {
	result    interface{}
	remaining parserInput
}

type search struct {
	results []*char
	desired interface{}
}

var (
	comments = []comment{
		comment{"//", "\n"},
		comment{"/*", "*/"},
	}
)

// DIGIT DIGIT DOT DIGIT DIGIT DOT....

// date -> number -> dot -> number -> dot -> year
// number -> digit digit
// dot
// number -> digit ..........

type parser func(parserInput) parserOutput

type commentInput struct {
	pos   int
	chars []*char // to allow nil
}

func (input commentInput) current() *char {
	if input.pos >= len(input.chars) {
		return nil
	}
	return input.chars[input.pos]
}

func (input commentInput) advance() parserInput {
	advancedPos := input.pos + 1
	if advancedPos >= len(input.chars) {
		return nil
	}
	return commentInput{advancedPos, input.chars}
}

func success(result interface{}, input parserInput) parserOutput {
	return parserOutput{result, input}
}

func fail(input parserInput) parserOutput {
	return parserOutput{nil, input}
}

// Get the size of a rune in UTF-16 format
func getUtf16RuneSize(r rune) int64 {
	rUtf16 := utf16.Encode([]rune{r}) // convert to utf16, since indices in GDocs API are utf16
	return int64(len(rUtf16))         // size of rune in utf16 format
}

// Get the size of a string in UTF-16 format
func getUtf16StringSize(s string) int64 {
	var size int64
	for _, r := range s {
		size += getUtf16RuneSize(r)
	}
	return size
}

// should i use ptr?
// string size must be > 0

// expectword /* expectspace or nothing, expectword */

type isRuneFunc func(r rune) bool

// if never gets what it is looking for, then whatever

// consume until the next thing is non-null?
func searchUntil(p parser) parser {
	return func(input parserInput) parserOutput {
		var results []*char
		output := p(input)
		for ; output.result == nil; output = p(input) {
			input = output.remaining
			output = expectChar(anyRune())(input)
			if output.result == nil {
				return success(search{results, nil}, input)
			}
			results = append(results, output.result.(*char))
			input = output.remaining
		}
		input = output.remaining
		return success(search{results, output.result}, input)
	}

}

func selectAny(parsers []parser) parser {
	return func(input parserInput) parserOutput {
		for _, p := range parsers {
			if output := p(input); output.result != nil {
				return output
			}
		}
		return fail(input)
	}
}

func separateComments(chars []*char) ([]*char, []*word) {
	var commentParsers []parser
	for _, c := range comments {
		commentParsers = append(commentParsers, expectComment(c.startSymbol, c.endSymbol))
	}
	var cs []*char
	var ws []*word
	var input parserInput = commentInput{0, chars}
	for input != nil && input.current() != nil {
		output := selectAny(commentParsers)(input)
		if output.result != nil {
			ws = append(ws, output.result.(*word)) // got a comment
			input = output.remaining
			continue
		}
		cs = append(cs, input.current())
		input = input.advance()
	}
	return cs, ws
}

func expectComment(start string, end string) parser {
	return func(input parserInput) parserOutput {
		output := expectWord(start)(input)
		if output.result == nil {
			return fail(input)
		}
		input = output.remaining
		w := output.result.(*word)
		var b bytes.Buffer
		b.WriteString(w.content)
		output = searchUntil(expectWord(end))(input)
		s := output.result.(search)
		for _, r := range s.results {
			w.index += r.size
			b.WriteRune(r.content)
		}
		if s.desired != nil {
			desired := s.desired.(*word)
			w.index += desired.size
			b.WriteString(desired.content)
		}
		input = output.remaining
		w.content = b.String()
		return success(w, input)
	}
}

func expectWord(s string) parser {
	return func(input parserInput) parserOutput {
		var w *word = nil
		for _, r := range s {
			output := expectChar(isRune(r))(input)
			if output.result == nil {
				return fail(input)
			}
			c := output.result.(*char)
			if w == nil {
				w = &word{c.index, 0, s}
			}
			w.size += c.size
			input = output.remaining
		}
		return success(w, input)
	}
}

func anyRune() isRuneFunc {
	return func(r1 rune) bool {
		return true
	}
}

func isRune(r1 rune) isRuneFunc {
	return func(r2 rune) bool {
		return r1 == r2
	}
}

func expectChar(desired isRuneFunc) parser {
	return func(input parserInput) parserOutput {
		if input == nil {
			return fail(input)
		}
		c := input.current()
		if c == nil || !desired(c.content) {
			return fail(input)
		}
		return success(c, input.advance())
	}
}

func getSlice(s string) []*char {
	var cs []*char
	for _, r := range s {
		cs = append(cs, &char{0, 0, r})
	}

	return cs
}

// replace comment with space
func main() {
	cs := getSlice("te// /*\nstinhig blah blah /* hello there */ what up")
	czz, wzz := separateComments(cs)
	for _, c := range czz {
		fmt.Printf("%c\n", c.content)
	}
	_ = czz
	_ = wzz
	// for _, w := range wzz {
	// 	fmt.Println(w)
	// }
	// w := output.result.(*word)
	// fmt.Printf("\n%v", w.content)
	// commInput = output.remaining.(commentInput)
	// fmt.Printf("\n\n%c\n", commInput.current().content)
	//fmt.Printf("\noutput: %s, remain: %v", output.result, output.remaining.(commentInput).chars[1:])
}
