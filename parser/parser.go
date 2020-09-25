package parser

import (
	"GDocs-Syntax-Highlighter/style"
	"fmt"
	"strings"
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
