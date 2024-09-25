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
)

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
	titleElem.SetText(m.Title)
	if m.TitleSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#title")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(m.TitleSort)
	}

	authElem := mdataElem.CreateElement("dc:creator")
	authElem.CreateAttr("id", "author")
	authElem.SetText(m.Author)
	if m.AuthorSort != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", "#author")
		metaElem.CreateAttr("property", "file-as")
		metaElem.SetText(m.AuthorSort)
	}

	mdataElem.CreateElement("dc:language").SetText(m.Language)

	if m.Series != "" {
		seriesElem := mdataElem.CreateElement("meta")
		seriesElem.CreateAttr("property", "belongs-to-collection")
		seriesElem.CreateAttr("id", "series")
		seriesElem.SetText(m.Series)

		if !math.IsNaN(m.SeriesNum) {
			seriesNumElem := mdataElem.CreateElement("meta")
			seriesNumElem.CreateAttr("refines", "#series")
			seriesNumElem.CreateAttr("property", "group-position")
			seriesNumElem.SetText(strconv.FormatFloat(m.SeriesNum, 'f', 2, 64))
		}
	}

	for _, s := range m.Subjects {
		mdataElem.CreateElement("dc:subject").SetText(s)
	}

	if m.Isbn != "" {
		isbnElem := mdataElem.CreateElement("dc:identifier")
		isbnElem.CreateAttr("opf:scheme", "ISBN")
		isbnElem.CreateAttr("id", "ISBN")
		isbnElem.SetText(m.Isbn)
	}

	if m.Publisher != "" {
		mdataElem.CreateElement("dc:publisher").SetText(m.Publisher)
	}

	if m.PubDate != "" {
		mdataElem.CreateElement("dc:date").SetText(m.PubDate)
	}

	if m.Rights != "" {
		mdataElem.CreateElement("dc:rights").SetText(m.Rights)
	}

	for i, c := range m.Contributors {
		id := fmt.Sprintf("contributor_%d", i)
		contElem := mdataElem.CreateElement("dc:contributor")
		contElem.CreateAttr("id", id)
		contElem.SetText(c.name)
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("refines", fmt.Sprintf("#%s", id))
		metaElem.CreateAttr("property", "role")
		metaElem.SetText(c.role)
	}

	if m.Description != "" {
		mdataElem.CreateElement("dc:description").SetText(m.Description)
	}

	idId := pkgElem.SelectAttrValue("unique-identifier", "uid")
	uidElem := mdataElem.CreateElement("dc:identifier")
	uidElem.CreateAttr("id", idId)
	uidElem.SetText(m.Uid)

	if coverId != "" {
		metaElem := mdataElem.CreateElement("meta")
		metaElem.CreateAttr("name", "cover")
		metaElem.CreateAttr("content", coverId)
	}
	pkgElem.InsertChildAt(0, mdataElem)

	return nil
}