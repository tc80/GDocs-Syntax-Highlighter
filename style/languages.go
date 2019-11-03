package style

// type Keyword struct {
// 	// Value        string
// 	KeywordColor Color
// 	//IsCaseSensitive bool
// }

type Comment struct {
	StartSymbol string
	EndSymbol   string
}

//BackgroundColor  Color - should be specified by user
//DefaultTextColor Color - should be specified by user

// Language represents a programming language
type Language struct {
	Keywords  map[string]Color
	Comments  []Comment
	Shortcuts map[string]string
	Font      string
}

var (
	languages = map[string]Language{
		"java": Language{
			Keywords: map[string]Color{
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
			Font: CourierNew,
		},
	}
)

// GetLanguage returns the Language specified
// by the string parameter if found, else nil
func GetLanguage(lang string) *Language {
	if l, ok := languages[lang]; ok {
		return &l
	}
	return nil
}
