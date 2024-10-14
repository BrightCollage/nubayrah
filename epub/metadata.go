package epub

type Metadata struct {
	Title        string        `json:"title"`
	TitleSort    string        `json:"titleSort"`
	Author       string        `json:"author"`
	AuthorSort   string        `json:"authorSort"`
	Language     string        `json:"language"`
	Series       string        `json:"series"`
	SeriesNum    float64       `json:"seriesNum"`
	Subjects     []string      `json:"subjects" gorm:"serializer:json" `
	Isbn         string        `json:"isbn"`
	Publisher    string        `json:"publisher"`
	PubDate      string        `json:"pubDate"` // iso8601 format
	Rights       string        `json:"rights"`
	Contributors []Contributor `json:"contributors" gorm:"serializer:json" `
	Description  string        `json:"description"`
	// The following fields are not user-editable
	Uid string `json:"uid"`
}

type Contributor struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
