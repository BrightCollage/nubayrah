package epub

import (
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
func ExtractMetadata(doc *RootFile) *Metadata {
	mdata := &Metadata{}

	mdata.title, mdata.titleSort = doc.getTitle()
	mdata.author, mdata.authorSort = doc.getAuthor()
	mdata.language = doc.getLanguage()
	mdata.series, mdata.seriesNum = doc.getSeries()
	mdata.subjects = doc.getSubjects()
	mdata.isbn = doc.getISBN()
	mdata.publisher = doc.getPublisher()
	mdata.pubDate = doc.getPubDate()
	mdata.rights = doc.getRights()
	mdata.contributors = doc.getContributors()
	mdata.description = doc.getDescription()

	mdata.uid = doc.getUID()
	mdata.nubayrahId = doc.getNubayrahId()

	return mdata
}

// Reads title, titleSort from xml doc
func (f *RootFile) getTitle() (title string, titleSort string) {

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
	seriesNum = math.NaN()
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
		contributors[i].name = c.Text()
		if r := c.SelectAttrValue("opf:role", ""); r != "" {
			contributors[i].role = r
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
			contributors[i].role = metaElem.Text()
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

func (f *RootFile) getNubayrahId() string {
	idElem := f.FindElement("//meta[@property='nubayrahId']")
	if idElem == nil {
		return ""
	}
	return idElem.Text()
}

func (f *RootFile) getCoverId() string {
	elem := f.FindElement("//meta[@name='cover']")
	if elem == nil {
		return ""
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
