package aws

import "time"

type Result struct {
	Status              string        `json:"status"`
	CreatedDateTime     time.Time     `json:"createdDateTime"`
	LastUpdatedDateTime time.Time     `json:"lastUpdatedDateTime"`
	AnalyzeResult       AnalyzeResult `json:"analyzeResult"`
}

type AnalyzeResult struct {
	APIVersion      string     `json:"apiVersion"`
	ModelID         string     `json:"modelId"`
	StringIndexType string     `json:"stringIndexType"`
	Content         string     `json:"content"`
	Pages           []Page     `json:"pages"`
	Tables          []Table    `json:"tables"`
	Styles          []Style    `json:"styles"`
	Documents       []Document `json:"documents"`
}

type Page struct {
	PageNumber int     `json:"pageNumber"`
	Angle      float64 `json:"angle"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Unit       string  `json:"unit"`
	Words      []Word  `json:"words"`
	Lines      []Line  `json:"lines"`
	Spans      []Span  `json:"spans"`
}

type Word struct {
	Content    string  `json:"content"`
	Polygon    []int   `json:"polygon"`
	Confidence float64 `json:"confidence"`
	Span       Span    `json:"span"`
}

type Line struct {
	Content string `json:"content"`
	Polygon []int  `json:"polygon"`
	Spans   []Span `json:"spans"`
}

type Span struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type Table struct {
	RowCount        int              `json:"rowCount"`
	ColumnCount     int              `json:"columnCount"`
	Cells           []TableCell      `json:"cells"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Spans           []Span           `json:"spans"`
}

type TableCell struct {
	Kind            string           `json:"kind,omitempty"`
	RowIndex        int              `json:"rowIndex"`
	ColumnIndex     int              `json:"columnIndex"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Spans           []Span           `json:"spans"`
	ColumnSpan      int              `json:"columnSpan,omitempty"`
	RowSpan         int              `json:"rowSpan,omitempty"`
}

type BoundingRegion struct {
	PageNumber int   `json:"pageNumber"`
	Polygon    []int `json:"polygon"`
}

type Style struct {
	Confidence    float64 `json:"confidence"`
	Spans         []Span  `json:"spans"`
	IsHandwritten bool    `json:"isHandwritten"`
}

type Document struct {
	DocType         string           `json:"docType"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Fields          DocumentFields   `json:"fields"`
	Confidence      float64          `json:"confidence"`
	Spans           []Span           `json:"spans"`
}

type DocumentFields struct {
	Birthday FieldDetail `json:"birthday"`
	Number   FieldDetail `json:"number"`
}

type FieldDetail struct {
	Type            string           `json:"type"`
	ValueString     string           `json:"valueString"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Confidence      float64          `json:"confidence"`
	Spans           []Span           `json:"spans"`
}
