package parser

import (
	"GDocs-Syntax-Highlighter/style"
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"google.golang.org/api/docs/v1"
)

const (
	// sometimes cannot find begin and end
	// need to fix
	beginSymbol = "~~begin~~"
	endSymbol   = "~~end~~"
)

// Char is a rune and its respective utf16 start and end indices
type Char struct {
	index   int64 // the utf16 inclusive start index of the rune
	size    int64 // the size of the rune in utf16 units
	content rune  // the rune
}

// Word is a string and its respective utf16 start and end indices
type Word struct {
	Index   int64  // the utf16 inclusive start index of the string
	Size    int64  // the size of the Word in utf16 units
	Content string // the string
}

type parser func(parserInput) parserOutput

// Comment
type comment struct {
	startSymbol string
	endSymbol   string
}

type parserInput interface {
	current() *Char
	advance() parserInput
}

type parserOutput struct {
	result    interface{}
	remaining parserInput
}

type search struct {
	results []*Char
	desired interface{}
}

type commentInput struct {
	pos   int
	chars []*Char // to allow nil
}

func (input commentInput) current() *Char {
	if input.pos >= len(input.chars) {
		return nil
	}
	return input.chars[input.pos]
}

func (input commentInput) advance() parserInput {
	advancedPos := input.pos + 1
	return commentInput{advancedPos, input.chars}
}

func success(result interface{}, input parserInput) parserOutput {
	return parserOutput{result, input}
}

func fail() parserOutput {
	return parserOutput{nil, nil}
}

// should i use ptr?
// string size must be > 0

// expectword /* expectspace or nothing, expectword */

// if never gets what it is looking for, then whatever

// consume until the next thing is non-null?
// strict means i must find the thing
func searchUntil(p parser) parser {
	return func(input parserInput) parserOutput {
		var results []*Char
		output := p(input)
		for ; output.result == nil; output = p(input) {
			output = expectChar(anyRune())(input)
			if output.result == nil {
				// reached end, did not find
				return success(search{results, nil}, input)
			}
			results = append(results, output.result.(*Char))
			input = output.remaining
		}
		// output.result will determine if we reached the end or not
		return success(search{results, output.result}, output.remaining)
	}

}

// Selects the first parser in a slice of
// parsers that successfully parses the input
func selectAny(parsers []parser) parser {
	return func(input parserInput) parserOutput {
		for _, p := range parsers {
			if output := p(input); output.result != nil {
				return output
			}
		}
		return fail()
	}
}

// Gets a filler character with size of 1
// and value as a space
func getFillerChar(index int64) *Char {
	return &Char{index, 1, ' '}
}

// SeparateComments does...
func SeparateComments(language style.Language, chars []*Char) ([]*Char, []*Word) {
	var commentParsers []parser
	for _, c := range language.Comments {
		commentParsers = append(commentParsers, expectComment(c.StartSymbol, c.EndSymbol))
	}
	var cs []*Char
	var ws []*Word
	var input parserInput = commentInput{0, chars}
	for input.current() != nil {
		output := selectAny(commentParsers)(input)
		if output.result != nil {
			w := output.result.(*Word)
			ws = append(ws, w)                      // got a comment
			cs = append(cs, getFillerChar(w.Index)) // put filler in for something like hello/**/world so it is hello world instead of helloworld
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
			return fail()
		}
		input = output.remaining
		w := output.result.(*Word)
		var b bytes.Buffer
		b.WriteString(w.Content)
		output = searchUntil(expectWord(end))(input)
		s := output.result.(search)
		for _, r := range s.results {
			w.Size += r.size
			b.WriteRune(r.content)
		}
		if s.desired != nil {
			fmt.Println("found")
			desired := s.desired.(*Word)
			w.Size += desired.Size
			b.WriteString(desired.Content)
		}
		input = output.remaining
		w.Content = b.String()
		fmt.Println(b.String())
		return success(w, input)
	}
}

func expectWord(s string) parser {
	return func(input parserInput) parserOutput {
		var w *Word
		for _, r := range s {
			output := expectChar(isRune(r))(input)
			if output.result == nil {
				return fail()
			}
			c := output.result.(*Char)
			if w == nil {
				w = &Word{c.index, 0, s}
			}
			w.Size += c.size
			input = output.remaining
		}
		return success(w, input)
	}
}

// Expects a given character based on
// a boolean character function
func expectChar(desired isRuneFunc) parser {
	return func(input parserInput) parserOutput {
		c := input.current()
		if c == nil || !desired(c.content) {
			return fail()
		}
		return success(c, input.advance())
	}
}

// GetSlice ...
func GetSlice(s string) []*Char {
	var cs []*Char
	for _, r := range s {
		cs = append(cs, &Char{0, 0, r})
	}

	return cs
}

// GetWords gets a Word slice from a Char slice
func GetWords(chars []*Char) []*Word {
	var words []*Word
	var b bytes.Buffer
	var index int64
	start := true
	for _, Char := range chars {
		if unicode.IsSpace(Char.content) {
			// we are separating words by space characters
			str := b.String()
			if len(str) > 0 {
				// word must have at least one char
				size := GetUtf16StringSize(str)
				words = append(words, &Word{index, size, str})
				start = true
				b.Reset()
			}
			continue
		}
		if start {
			index = Char.index
			start = false
		}
		b.WriteRune(Char.content)
	}
	str := b.String()
	if len(str) > 0 {
		size := GetUtf16StringSize(str)
		words = append(words, &Word{index, size, str})
	}
	return words
}

// GetChars gets the slice of all chars, where
// each Char holds a rune and its respective utf16 range
func GetChars(doc *docs.Document) []*Char {
	var chars []*Char
	begin := false
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, par := range elem.Paragraph.Elements {
				if par.TextRun != nil {
					content := strings.TrimSpace(par.TextRun.Content)
					fmt.Println(content)
					if strings.EqualFold(content, endSymbol) {
						return chars
					}
					if !begin {
						if strings.EqualFold(content, beginSymbol) {
							begin = true
						}
						continue
					}
					index := par.StartIndex
					// iterate over runes
					for _, r := range par.TextRun.Content {
						size := GetUtf16RuneSize(r)                  // size of run in utf16 units
						chars = append(chars, &Char{index, size, r}) // associate runes with ranges
						index += size
					}
				}
			}
		}
	}
	return chars
}

// GetRange ...
func GetRange(chars []*Char) (int64, int64) {
	startIndex := chars[0].index
	lastChar := chars[len(chars)-1]
	endIndex := lastChar.index + lastChar.size
	return startIndex, endIndex
}
