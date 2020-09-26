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
	underlineField     = "underline"
	smallCapsField     = "smallCaps"
	strikethroughField = "strikethrough"
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

// SetUnderline sets a range to underline or not.
func SetUnderline(underline bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: underlineField,
			Range:  r,
			TextStyle: &docs.TextStyle{
				Underline: underline,
			},
		},
	}
}

// SetSmallCaps sets a range to small caps or not.
func SetSmallCaps(smallCaps bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: smallCapsField,
			Range:  r,
			TextStyle: &docs.TextStyle{
				SmallCaps: smallCaps,
			},
		},
	}
}

// SetStrikethrough sets a range to strikethrough or not.
func SetStrikethrough(strikethrough bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: strikethroughField,
			Range:  r,
			TextStyle: &docs.TextStyle{
				Strikethrough: strikethrough,
			},
		},
	}
}

// ClearFormatting removes any italics, bold, smallcaps, strikethrough,
// and underline.
func ClearFormatting(r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields:    getFields(italicField, boldField, smallCapsField, strikethroughField, underlineField),
			Range:     r,
			TextStyle: &docs.TextStyle{},
		},
	}
}

// UpdateFont gets the request to update a range with a particular font.
// Note, this will unbold this range due to weightedFontFamily.
func UpdateFont(font string, size float64, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: getFields(weightedFontFamily, fontSize),
			Range:  r,
			TextStyle: &docs.TextStyle{
				WeightedFontFamily: &docs.WeightedFontFamily{
					FontFamily: font,
				},
				FontSize: &docs.Dimension{
					Magnitude: size,
					Unit:      pointUnit,
				},
			},
		},
	}
}

// UpdateTextStyleExceptUnderline updates a font and formatting options except for underline,
// since underline is used in the directives.
func UpdateTextStyleExceptUnderline(font string, size float64, italic, bold, smallCaps, strikethrough bool, r *docs.Range) *docs.Request {
	return &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Fields: getFields(weightedFontFamily, fontSize, italicField, boldField, smallCapsField, strikethroughField),
			Range:  r,
			TextStyle: &docs.TextStyle{
				WeightedFontFamily: &docs.WeightedFontFamily{
					FontFamily: font,
				},
				FontSize: &docs.Dimension{
					Magnitude: size,
					Unit:      pointUnit,
				},
				Italic:        italic,
				Bold:          bold,
				SmallCaps:     smallCaps,
				Strikethrough: strikethrough,
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
