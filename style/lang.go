package style

import "strings"

// type Keyword struct {
// 	// Value        string
// 	KeywordColor Color
// 	//IsCaseSensitive bool
// }

// Comment ...
type Comment struct {
	StartSymbol string
	EndSymbol   string
}

//BackgroundColor  Color - should be specified by user
//DefaultTextColor Color - should be specified by user

// Language represents a programming language
type Language struct {
	Name      string
	Keywords  map[string]*Color
	Comments  []Comment
	Shortcuts map[string]string
	Font      string
}

var (
	javaLang = Language{
		Name: "Java",
		Keywords: map[string]*Color{
			"public": Red,
			"static": Blue,
			"void":   Green,
			"if":     Blue,
		},
		Comments: []Comment{
			{"//", "\n"},
			{"/*", "*/"},
		},
		Shortcuts: map[string]string{
			"psvm":  "public static void main(String[] args) {\n\n}",
			"if-el": "if (cond) {\n\n} else {\n\n}",
		},
		Font: CourierNew, // maybe user defines font ?
	}
	goLang = Language{
		Name: "Go",
		Keywords: map[string]*Color{
			"package":     Red,
			"func":        Blue,
			"\"":          Green,
			"fmt.Println": Blue,
		},
		Comments: []Comment{
			{"//", "\n"},
			{"/*", "*/"},
		},
		Shortcuts: map[string]string{
			"psvm":  "public static void main(String[] args) {\n\n}",
			"if-el": "if (cond) {\n\n} else {\n\n}",
		},
		Font: CourierNew, // maybe user defines font ?
	}
	languages = map[string]Language{
		"java": javaLang,
		"go":   goLang,
	}
)

// GetLanguage attempts to get a Language
// from a case insensitive string.
func GetLanguage(lang string) (Language, bool) {
	l, ok := languages[strings.ToLower(lang)]
	return l, ok
}

// GetDefaultLanguage gets the default Language
// if the directive is not set.
func GetDefaultLanguage() Language {
	return goLang
}
