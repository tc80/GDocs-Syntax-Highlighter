package requests

import (
	"GDocs-Syntax-Highlighter/parser"
	"GDocs-Syntax-Highlighter/style"

	"google.golang.org/api/docs/v1"
)

const (
	background         = "background"
	foregroundColor    = "foregroundColor"
	backgroundColor    = "backgroundColor"
	weightedFontFamily = "weightedFontFamily"
)

// GetDocumentColorRequest gets a request to change the color of the document.
func GetDocumentColorRequest(c *style.Color) *docs.Request {
	return &docs.Request{
		UpdateDocumentStyle: &docs.UpdateDocumentStyleRequest{
			Fields: background,
			DocumentStyle: &docs.DocumentStyle{
				Background: &docs.Background{
					Color: &docs.OptionalColor{
						Color: &docs.Color{
							RgbColor: &docs.RgbColor{
								Blue:  c.Blue,
								Red:   c.Red,
								Green: c.Green,
							},
						},
					},
				},
			},
		},
	}
}

// GetForeColorRequest gets a request to change the foreground color of a range.
func GetForeColorRequest(c *style.Color, startIndex, endIndex int64) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: foregroundColor,
			Range: &docs.Range{
				StartIndex: startIndex,
				EndIndex:   endIndex,
			},
			TextStyle: &docs.TextStyle{
				ForegroundColor: &docs.OptionalColor{
					Color: &docs.Color{
						RgbColor: &docs.RgbColor{
							Red:   c.Red,
							Blue:  c.Blue,
							Green: c.Green,
						},
					},
				},
			},
		},
	}
}

// GetBackColorRequest gets a request to change the background color of a range.
func GetBackColorRequest(c *style.Color, startIndex, endIndex int64) *docs.Request {
	var color *docs.Color
	// if c is nil, it is transparent
	if c != nil {
		color = &docs.Color{
			RgbColor: &docs.RgbColor{
				Red:   c.Red,
				Blue:  c.Blue,
				Green: c.Green,
			},
		}
	}
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: backgroundColor,
			Range: &docs.Range{
				StartIndex: startIndex,
				EndIndex:   endIndex,
			},
			TextStyle: &docs.TextStyle{
				BackgroundColor: &docs.OptionalColor{
					Color: color,
				},
			},
		},
	}
}

// GetDeleteRequest ...
func GetDeleteRequest(start, end int64) *docs.Request {
	return &docs.Request{
		DeleteContentRange: &docs.DeleteContentRangeRequest{
			Range: &docs.Range{
				StartIndex: start,
				EndIndex:   end,
			},
		},
	}
}

// GetInsertRequest ...
func GetInsertRequest(text string, start int64) *docs.Request {
	return &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Text: text,
			Location: &docs.Location{
				Index: start,
			},
		},
	}
}

// GetReplaceRequest gets the requests to delete a Word and insert a new one in its place.
func GetReplaceRequest(word *parser.Word, wordsAfter []*parser.Word, replace string) []*docs.Request {
	// request to delete the Word
	delete := GetDeleteRequest(word.Index, word.Index+word.Size)

	// request to insert the replacement at deleted Word's location
	insert := GetInsertRequest(replace, word.Index)

	requests := []*docs.Request{delete, insert}
	newSize := parser.GetUtf16StringSize(replace)
	diff := newSize - word.Size
	word.Size = newSize
	// update ranges for Words that follow this Word
	for _, w := range wordsAfter {
		w.Index += diff
	}
	return requests
}

// GetFontRequest gets the request to update a range with a particular font.
func GetFontRequest(font string, startIndex, endIndex int64) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: weightedFontFamily,
			Range: &docs.Range{
				StartIndex: startIndex,
				EndIndex:   endIndex,
			},
			TextStyle: &docs.TextStyle{
				WeightedFontFamily: &docs.WeightedFontFamily{
					FontFamily: font,
				},
			},
		},
	}
}

// GetBatchUpdate gets the batch request from a slice of requests.
func GetBatchUpdate(requests []*docs.Request) *docs.BatchUpdateDocumentRequest {
	return &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
}
