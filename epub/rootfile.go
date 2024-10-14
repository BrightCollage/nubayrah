/*
https://idpf.org/epub/20/spec/OPF_2.0.1_draft.htm
https://www.w3.org/TR/epub-33/
*/

package epub

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

type RootFile struct {
	*etree.Document
	internalPath string
}

// Parses metadata from root-file xml document
func (e *Epub) ExtractMetadata() *Metadata {
	mdata := &Metadata{}

	mdata.Title, mdata.TitleSort = e.RootFile.getTitle()
	mdata.Author, mdata.AuthorSort = e.RootFile.getAuthor()
	mdata.Language = e.RootFile.getLanguage()
	mdata.Series, mdata.SeriesNum = e.RootFile.getSeries()
	mdata.Subjects = e.RootFile.getSubjects()
	mdata.Isbn = e.RootFile.getISBN()
	mdata.Publisher = e.RootFile.getPublisher()
	mdata.PubDate = e.RootFile.getPubDate()
	mdata.Rights = e.RootFile.getRights()
	mdata.Contributors = e.RootFile.getContributors()
	mdata.Description = e.RootFile.getDescription()

	mdata.Uid = e.RootFile.getUID()

	return mdata
}

// Generates <metadata> element and inserts it as first child of <package>,
// overwriting any existing metadata element
func (f *RootFile) InsertMetadata(mdata *Metadata) error {
	pkgElem := f.FindElement("//package")
	if pkgElem == nil {
		return errors.New("malformed package document: package element not found")
	}

	coverId := f.getCoverId()

	for _, child := range pkgElem.ChildElements() {
		if child.Tag == "metadata" {
			pkgElem.RemoveChild(child)
			break
		}
	}

	mdataElem := pkgElem.CreateElement("metadata")
	pkgElem.RemoveChild(mdataElem)

	titleElem := mdataElem.CreateElement("dc:title")
	titleElem.CreateAttr("id", "title")
	titleElem.SetText(mdata.Title)
	if mdata.TitleSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#title")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(mdata.TitleSort)
	}

	authElem := mdataElem.CreateElement("dc:creator")
	authElem.CreateAttr("id", "author")
	authElem.SetText(mdata.Author)
	if mdata.AuthorSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#author")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(mdata.AuthorSort)
	}

	mdataElem.CreateElement("dc:language").SetText(mdata.Language)

	if mdata.Series != "" {
		seriesElem := mdataElem.CreateElement("meta")
		seriesElem.CreateAttr("property", "belongs-to-collection")
		seriesElem.CreateAttr("id", "series")
		seriesElem.SetText(mdata.Series)

		if !math.IsNaN(mdata.SeriesNum) {
			seriesNumElem := mdataElem.CreateElement("meta")
			seriesNumElem.CreateAttr("refines", "#series")
			seriesNumElem.CreateAttr("property", "group-position")
			seriesNumElem.SetText(strconv.FormatFloat(mdata.SeriesNum, 'f', 2, 64))
		}
	}

	for _, s := range mdata.Subjects {
		mdataElem.CreateElement("dc:subject").SetText(s)
	}

	if mdata.Isbn != "" {
		isbnElem := mdataElem.CreateElement("dc:identifier")
		isbnElem.CreateAttr("opf:scheme", "ISBN")
		isbnElem.CreateAttr("id", "ISBN")
		isbnElem.SetText(mdata.Isbn)
	}

	if mdata.Publisher != "" {
		mdataElem.CreateElement("dc:publisher").SetText(mdata.Publisher)
	}

	if mdata.PubDate != "" {
		mdataElem.CreateElement("dc:date").SetText(mdata.PubDate)
	}

	if mdata.Rights != "" {
		mdataElem.CreateElement("dc:rights").SetText(mdata.Rights)
	}

	for i, c := range mdata.Contributors {
		id := fmt.Sprintf("contributor_%d", i)
		contElem := mdataElem.CreateElement("dc:contributor")
		contElem.CreateAttr("id", id)
		contElem.SetText(c.Name)
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", fmt.Sprintf("#%s", id))
		metaElem.CreateAttr("property", "role")
		metaElem.SetText(c.Role)
	}

	if mdata.Description != "" {
		mdataElem.CreateElement("dc:description").SetText(mdata.Description)
	}

	idId := pkgElem.SelectAttrValue("unique-identifier", "uid")
	uidElem := mdataElem.CreateElement("dc:identifier")
	uidElem.CreateAttr("id", idId)
	uidElem.SetText(mdata.Uid)

	if coverId != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("name", "cover")
		metaElem.CreateAttr("content", coverId)
	}
	pkgElem.InsertChildAt(0, mdataElem)

	return nil
}

// Reads title, titleSort from xml doc
func (f *RootFile) getTitle() (title string, titleSort string) {
	title = "Unknown Title"
	titleElem := f.FindElement("//dc:title")
	if titleElem == nil {
		return
	}

	title = titleElem.Text()
	titleSort = titleElem.SelectAttrValue("opf:file-as", "")
	if titleSort != "" {
		return
	}

	id := titleElem.SelectAttrValue("id", "")
	if id == "" {
		return
	}
	metaElem := f.FindElementFiltered("//*", filter{name: "property", value: "file-as"}, filter{name: "refines", value: fmt.Sprintf("#%s", id)})
	if metaElem != nil {
		titleSort = metaElem.Text()
		return
	}
	return
}

// Reads author and authorSort from xml doc
func (f *RootFile) getAuthor() (author string, authorSort string) {
	author = "Unknown Author"
	metaElem := f.FindElement("//dc:creator")
	if metaElem == nil {
		return
	}
	author = metaElem.Text()
	authorSort = metaElem.SelectAttrValue("opf:file-as", "")
	if authorSort != "" {
		return
	}

	id := metaElem.SelectAttrValue("id", "")
	if id == "" {
		return
	}
	metaElem = f.FindElementFiltered(
		"//meta",
		filter{
			name:  "property",
			value: "file-as",
		},
		filter{
			name:  "refines",
			value: fmt.Sprintf("#%s", id),
		})
	if metaElem != nil {
		authorSort = metaElem.Text()
	}
	return
}

// Reads series, seriesNum from xml doc
// If there is no seriesNum the value will be -1
func (f *RootFile) getSeries() (series string, seriesNum float64) {
	seriesNum = -1
	seriesMetaElem := f.FindElement("//meta[@property='belongs-to-collection']")
	if seriesMetaElem != nil {
		series = seriesMetaElem.Text()
		elemId := seriesMetaElem.SelectAttrValue("id", "")
		if elemId != "" {
			seriesNumMetaElem := f.FindElementFiltered(
				"//meta",
				filter{
					name:  "refines",
					value: fmt.Sprintf("#%s", elemId),
				},
				filter{
					name:  "property",
					value: "group-position",
				})
			if seriesNumMetaElem != nil {
				num, err := strconv.ParseFloat(seriesNumMetaElem.Text(), 64)
				if err == nil {
					seriesNum = num
				}
			}
		}
	}
	return
}

func (f *RootFile) getSubjects() []string {
	subjectElem := f.FindElements("//subject")
	subjects := make([]string, len(subjectElem))
	for i, s := range subjectElem {
		subjects[i] = s.Text()
	}
	return subjects
}

func (f *RootFile) getLanguage() string {
	return f.getNodeText("dc:language")
}

func (f *RootFile) getISBN() string {
	isbnElem := f.FindElement("//*[@opf:scheme='ISBN']")
	if isbnElem == nil {
		isbnElem = f.FindElement("//*[@opf:scheme='isbn']")
	}

	if isbnElem != nil {
		return isbnElem.Text()
	}

	return ""
}

func (f *RootFile) getPublisher() string {
	return f.getNodeText("dc:publisher")
}

func (f *RootFile) getPubDate() string {
	date := f.getNodeText("dc:date")
	return strings.Split(date, "T")[0]
}

func (f *RootFile) getRights() string {
	return f.getNodeText("dc:rights")
}

func (f *RootFile) getDescription() string {
	return f.getNodeText("description")
}

func (f *RootFile) getContributors() []Contributor {
	contribs := f.FindElements("//contributor")
	contributors := make([]Contributor, len(contribs))
	for i, c := range contribs {
		contributors[i].Name = c.Text()
		if r := c.SelectAttrValue("opf:role", ""); r != "" {
			contributors[i].Role = r
			continue
		}
		id := c.SelectAttrValue("id", "")
		if id == "" {
			continue
		}

		metaElem := f.FindElementFiltered(
			"//meta",
			filter{
				name:  "property",
				value: "role"},
			filter{
				name:  "refines",
				value: fmt.Sprintf("#%s", id),
			})

		if metaElem != nil {
			contributors[i].Role = metaElem.Text()
		}

	}

	return contributors
}

// Gets uid and id name of uid element indicated by package's unique-identifier
func (f *RootFile) getUID() string {
	pkgElem := f.FindElement("//package")
	if pkgElem == nil {
		return ""
	}

	uidName := pkgElem.SelectAttrValue("unique-identifier", "")
	if uidName == "" {
		return ""
	}

	uidElem := f.FindElement(fmt.Sprintf("//dc:identifier[@id='%s']", uidName))

	if uidElem == nil {
		return ""
	}

	return uidElem.Text()
}

func (f *RootFile) getCoverId() string {
	elem := f.FindElement("//meta[@name='cover']")
	if elem == nil {
		return "cover-image"
	}
	return elem.SelectAttrValue("content", "")
}

// Gets text value of node or an empty string
func (f *RootFile) getNodeText(name string) string {
	elem := f.FindElement(fmt.Sprintf("//%s", name))
	if elem == nil {
		return ""
	}
	return elem.Text()
}

/*
Beevik's implementation of etree doesn't work well with multiple filters eg:
//meta[@color='blue' and @shape='square']

This method will find the first element that satisfies the provided filters,
though it will be slightly slower so use FindElement(xpath) if you have only
one filter.
*/
func (f *RootFile) FindElementFiltered(xpath string, filters ...filter) *etree.Element {
	for _, elem := range f.FindElements(xpath) {
		allOk := true
		for _, filter := range filters {
			if elem.SelectAttrValue(filter.name, "") != filter.value {
				allOk = false
				break
			}
		}
		if allOk {
			return elem
		}
	}
	return nil
}

type filter struct {
	name  string
	value string
}
