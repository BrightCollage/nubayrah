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

	"github.com/google/uuid"
)

type Metadata struct {
	title        string
	titleSort    string
	author       string
	authorSort   string
	language     string
	series       string
	seriesNum    float64
	subjects     []string
	isbn         string
	publisher    string
	pubDate      string // iso8601 format
	rights       string
	contributors []Contributor
	description  string
	// The following fields are not user-editable
	uid        string
	nubayrahId string
}

type Contributor struct {
	name string
	role string
}

func (c *Contributor) ToString() string {
	return fmt.Sprintf("%s: %s", c.name, c.role)
}

// Generates <metadata> element and inserts it as first child of <package>,
// overwriting any existing metadata element
func (m *Metadata) RenderToDoc(rootFile *RootFile) error {

	pkgElem := rootFile.FindElement("//package")
	if pkgElem == nil {
		return errors.New("malformed package document: package element not found")
	}

	coverId := rootFile.getCoverId()

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
	titleElem.SetText(m.title)
	if m.titleSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#title")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(m.titleSort)
	}

	authElem := mdataElem.CreateElement("dc:creator")
	authElem.CreateAttr("id", "author")
	authElem.SetText(m.author)
	if m.authorSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#author")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(m.authorSort)
	}

	mdataElem.CreateElement("dc:language").SetText(m.language)

	if m.series != "" {
		seriesElem := mdataElem.CreateElement("meta")
		seriesElem.CreateAttr("property", "belongs-to-collection")
		seriesElem.CreateAttr("id", "series")
		seriesElem.SetText(m.series)

		if !math.IsNaN(m.seriesNum) {
			seriesNumElem := mdataElem.CreateElement("meta")
			seriesNumElem.CreateAttr("refines", "#series")
			seriesNumElem.CreateAttr("property", "group-position")
			seriesNumElem.SetText(strconv.FormatFloat(m.seriesNum, 'f', 2, 64))
		}
	}

	for _, s := range m.subjects {
		mdataElem.CreateElement("dc:subject").SetText(s)
	}

	if m.isbn != "" {
		isbnElem := mdataElem.CreateElement("dc:identifier")
		isbnElem.CreateAttr("opf:scheme", "ISBN")
		isbnElem.CreateAttr("id", "ISBN")
		isbnElem.SetText(m.isbn)
	}

	if m.publisher != "" {
		mdataElem.CreateElement("dc:publisher").SetText(m.publisher)
	}

	if m.pubDate != "" {
		mdataElem.CreateElement("dc:date").SetText(m.pubDate)
	}

	if m.rights != "" {
		mdataElem.CreateElement("dc:rights").SetText(m.rights)
	}

	for i, c := range m.contributors {
		id := fmt.Sprintf("contributor_%d", i)
		contElem := mdataElem.CreateElement("dc:contributor")
		contElem.CreateAttr("id", id)
		contElem.SetText(c.name)
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", fmt.Sprintf("#%s", id))
		metaElem.CreateAttr("property", "role")
		metaElem.SetText(c.role)
	}

	if m.description != "" {
		mdataElem.CreateElement("dc:description").SetText(m.description)
	}

	idId := pkgElem.SelectAttrValue("unique-identifier", "uid")
	uidElem := mdataElem.CreateElement("dc:identifier")
	uidElem.CreateAttr("id", idId)
	uidElem.SetText(m.uid)

	if m.nubayrahId == "" {
		m.nubayrahId = uuid.New().String()
	}
	nidElem := mdataElem.CreateElement("meta")
	nidElem.CreateAttr("property", "nubayrahId")
	nidElem.SetText(m.nubayrahId)

	if coverId != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("name", "cover")
		metaElem.CreateAttr("content", coverId)
	}
	pkgElem.InsertChildAt(0, mdataElem)

	return nil
}
