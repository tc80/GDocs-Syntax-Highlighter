package request

import (
	"google.golang.org/api/docs/v1"
)

const (
	shading            = "shading"
	background         = "background"
	foregroundColor    = "foregroundColor"
	backgroundColor    = "backgroundColor"
	weightedFontFamily = "weightedFontFamily"
	fontSize           = "fontSize"
	boldField          = "bold"
	italicField        = "italic"
	pointUnit          = "PT"
	startIndex         = "StartIndex"
	endIndex           = "EndIndex"
)

// UpdateDocBackground gets a request to change the background color of the document.
func UpdateDocBackground(c *docs.Color) *docs.Request {
	return &docs.Request{
		UpdateDocumentStyle: &docs.UpdateDocumentStyleRequest{
			Fields: background,
			DocumentStyle: &docs.DocumentStyle{
				Background: &docs.Background{
					Color: &docs.OptionalColor{
						Color: c,
					},
				},
			},
		},
	}
}

// UpdateForegroundColor gets a request to change the foreground color of a range.
func UpdateForegroundColor(c *docs.Color, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: foregroundColor,
			Range:  r,
			TextStyle: &docs.TextStyle{
				ForegroundColor: &docs.OptionalColor{
					Color: c,
				},
			},
		},
	}
}

// UpdateHighlightColor gets a request to change the highlight color of a range.
func UpdateHighlightColor(c *docs.Color, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: backgroundColor,
			Range:  r,
			TextStyle: &docs.TextStyle{
				BackgroundColor: &docs.OptionalColor{
					Color: c,
				},
			},
		},
	}
}

// UpdateBackgroundColor gets a request to change the background color of a range.
func UpdateBackgroundColor(c *docs.Color, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
			Fields: shading,
			Range:  r,
			ParagraphStyle: &docs.ParagraphStyle{
				Shading: &docs.Shading{
					BackgroundColor: &docs.OptionalColor{
						Color: c,
					},
				},
			},
		},
	}
}

// Insert inserts text at an index.
func Insert(text string, start int64) *docs.Request {
	return &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Text: text,
			Location: &docs.Location{
				Index: start,
			},
		},
	}
}

// Delete removes text in a range.
func Delete(r *docs.Range) *docs.Request {
	return &docs.Request{
		DeleteContentRange: &docs.DeleteContentRangeRequest{
			Range: r,
		},
	}
}

// SetBold sets a range to bold or not bold.
func SetBold(bold bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: boldField,
			Range:  r,
			TextStyle: &docs.TextStyle{
				Bold: bold,
			},
		},
	}
}

// SetItalics sets a range to italic or not italic.
func SetItalics(italic bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: italicField,
			Range:  r,
			TextStyle: &docs.TextStyle{
				Italic: italic,
			},
		},
	}
}

// UpdateFont gets the request to update a range with a particular font.
func UpdateFont(font string, size float64, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: getFields(weightedFontFamily, fontSize),
			Range:  r,
			TextStyle: &docs.TextStyle{
				FontSize: &docs.Dimension{
					Magnitude: size,
					Unit:      pointUnit,
				},
				WeightedFontFamily: &docs.WeightedFontFamily{
					FontFamily: font,
				},
			},
		},
	}
}

// GetRange gets a new *docs.Range for
// start and end indices.
func GetRange(start, end int64, segmentID string) *docs.Range {
	return &docs.Range{
		StartIndex: start,
		EndIndex:   end,
		SegmentId:  segmentID,
		// force send since a value of 0 in a header/footer
		// will be omitted in the JSON, causing a bad request
		ForceSendFields: []string{startIndex, endIndex},
	}
}

// BatchUpdate gets the batch request from a slice of requests.
func BatchUpdate(requests []*docs.Request) *docs.BatchUpdateDocumentRequest {
	return &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
}
