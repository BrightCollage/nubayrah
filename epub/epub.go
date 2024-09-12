package epub

import (
	"archive/zip"
	"errors"
	"fmt"

	"github.com/beevik/etree"
	"golang.org/x/net/html/charset"
)

type Epub struct {
	Filepath string
	Metadata *Metadata
}

func NewEpubFromFile(filepath string) (*Epub, error) {
	epub := &Epub{}
	epub.Filepath = filepath

	if err := epub.ReadMetadata(); err != nil {
		return nil, err
	}

	return epub, nil
}

func (e *Epub) ReadMetadata() error {

	rf, err := e.GetRootFile()
	if err != nil {
		return err
	}
	e.Metadata = NewMetadataFromXML(rf)
	return nil
}

// Reads and parses rootfile from epub
func (e *Epub) GetRootFile() (*RootFile, error) {
	containerFile, err := e.ReadFile("META-INF/container.xml")
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	doc.ReadSettings.CharsetReader = charset.NewReaderLabel
	err = doc.ReadFromBytes(containerFile)

	if err != nil {
		return nil, err
	}

	elemrf := doc.FindElement("//rootfile")
	if elemrf == nil {
		return nil, errors.New("invalid container xml: rootfile element not found")
	}

	for _, a := range elemrf.Attr {
		if a.Key == "full-path" {
			b, err := e.ReadFile(a.Value)
			if err != nil {
				return nil, err
			}

			doc := etree.NewDocument()
			doc.ReadSettings.CharsetReader = charset.NewReaderLabel
			doc.ReadSettings.ValidateInput = false
			err = doc.ReadFromBytes(b)

			return &RootFile{*doc}, err
		}
	}

	return nil, errors.New("invalid container xml: rootfile path not found")
}

// Reads file from zip into byte array
func (e *Epub) ReadFile(absPath string) ([]byte, error) {
	epubZip, err := zip.OpenReader(e.Filepath)
	if err != nil {
		return nil, err
	}
	defer epubZip.Close()

	var file *zip.File

	for _, f := range epubZip.File {
		if f.Name == absPath {
			file = f
			break
		}
	}

	if file == nil {
		return nil, fmt.Errorf("file not found: %s", absPath)
	}

	c, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	buf := make([]byte, int(file.UncompressedSize64))
	c.Read(buf)

	return buf, nil
}
