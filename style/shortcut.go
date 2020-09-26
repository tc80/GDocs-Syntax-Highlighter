package style

import (
	"regexp"
)

const (
	// DefaultShortcutSetting is whether shortcuts are enabled by
	// default.
	DefaultShortcutSetting = true
)

// Shortcut represents a shortcut.
// Part of preprocessing, regex matches are replaced by respective strings.
type Shortcut struct {
	Regex   *regexp.Regexp
	Replace string
}

var (
	doubleQuotes   = &Shortcut{regexp.MustCompile("“|”"), "\""}
	singleQuotes   = &Shortcut{regexp.MustCompile("‘|’"), "'"}
	goMainShortcut = &Shortcut{
		regexp.MustCompile("\\*\\*main\\*\\*"),
		"package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"hello world\")\n}\n",
	}
)
