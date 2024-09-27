package epub

type Metadata struct {
	Title        string
	TitleSort    string
	Author       string
	AuthorSort   string
	Language     string
	Series       string
	SeriesNum    float64
	Subjects     []string
	Isbn         string
	Publisher    string
	PubDate      string // iso8601 format
	Rights       string
	Contributors []Contributor
	Description  string
	// The following fields are not user-editable
	Uid string
}

type Contributor struct {
	name string
	role string
}
