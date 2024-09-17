package epub

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/go-xmlfmt/xmlfmt"
	"golang.org/x/net/html/charset"
)

type Epub struct {
	filepath   string
	Metadata   *Metadata
	fileHandle *zip.ReadCloser
	rootFile   *RootFile
	coverImage []byte
}

func OpenEpub(filepath string) (*Epub, error) {
	epub := &Epub{}
	epub.filepath = filepath

	return epub, epub.Reload()
}

// Reopens and reads epub from disk
func (e *Epub) Reload() error {
	e.Close()
	var err error = nil
	e.fileHandle, err = zip.OpenReader(e.filepath)
	if err != nil {
		return err
	}

	e.rootFile, err = e.getRootFile()
	if err != nil {
		return err
	}

	if err := e.readMetadata(); err != nil {
		return err
	}

	return nil
}

// Closes open file handles and renders epub invalid for writing
func (e *Epub) Close() {
	if e.fileHandle != nil {
		e.fileHandle.Close()
		e.fileHandle = nil
	}
	e.Metadata = nil
}

// Parse metadata from rootfile into e.Metadata
func (e *Epub) readMetadata() error {
	rf, err := e.getRootFile()
	if err != nil {
		return err
	}
	e.Metadata = ExtractMetadata(rf)
	return nil
}

// Reads and parses rootfile from epub
func (e *Epub) getRootFile() (*RootFile, error) {
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

	rootfileElem := doc.FindElement("//rootfile")
	if rootfileElem == nil {
		return nil, errors.New("invalid container xml: rootfile element not found")
	}

	rootfilePath := rootfileElem.SelectAttrValue("full-path", "")
	if rootfilePath == "" {
		return nil, errors.New("invalid container xml: rootfile element missing full-path attr")
	}

	b, err := e.ReadFile(rootfilePath)
	if err != nil {
		return nil, err
	}

	doc = etree.NewDocument()
	doc.ReadSettings.CharsetReader = charset.NewReaderLabel
	doc.ReadSettings.ValidateInput = false
	err = doc.ReadFromBytes(b)

	return &RootFile{doc, rootfilePath}, err
}

// Reads file from zip into byte array
func (e *Epub) ReadFile(path string) ([]byte, error) {
	var file *zip.File
	path = filepath.FromSlash(path)
	for _, f := range e.fileHandle.File {
		n := filepath.FromSlash(f.Name)
		if n == path {
			file = f
			break
		}
	}

	if file == nil {
		return nil, fmt.Errorf("file not found: %s", path)
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

// Unpacks epub to destination directory
func (e *Epub) unpack(destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	for _, file := range e.fileHandle.File {
		outFilePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(outFilePath, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		err := os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm)
		if err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

// Writes changes to metadata and cover image to epub file
func (e *Epub) WriteChanges() error {
	tmpDir := filepath.Join(os.TempDir(), filepath.Base(strings.Split(e.filepath, ".")[0]))
	tmpFile := tmpDir + ".epub"

	e.unpack(tmpDir)

	defer os.Remove(tmpDir)

	// RootFile + Metadata
	e.Metadata.RenderToDoc(e.rootFile)
	rfStr, err := e.rootFile.WriteToString()
	if err != nil {
		return err
	}
	rfStr = xmlfmt.FormatXML(rfStr, "", "  ")

	err = os.WriteFile(filepath.Join(tmpDir, e.rootFile.internalPath), []byte(rfStr), os.ModePerm)
	if err != nil {
		return err
	}

	// Cover image
	if len(e.coverImage) != 0 {
		coverPath, err := e.getCoverPath()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(tmpDir, coverPath), e.coverImage, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	fsys := os.DirFS(tmpDir)

	err = zipWriter.AddFS(fsys)
	if err != nil {
		return err
	}

	zipWriter.Close()
	file.Close()
	e.fileHandle.Close()
	err = os.Rename(tmpFile, e.filepath)
	if err != nil {
		return err
	}

	return nil
}

// Returns internal path to cover image
func (e *Epub) getCoverPath() (string, error) {
	metaElem := e.rootFile.FindElement("//package/metadata/meta[@name='cover']")
	if metaElem == nil {
		return "", errors.New("cover id not found in metadata")
	}

	coverId := metaElem.SelectAttrValue("content", "")
	if coverId == "" {
		return "", errors.New("cover meta element missing content attribute")
	}

	itemElem := e.rootFile.FindElement(fmt.Sprintf("//package/manifest/item[@id='%s']", coverId))
	if itemElem == nil {
		return "", errors.New("cover image item not found in manifest")
	}

	imgRelativePath := itemElem.SelectAttrValue("href", "")

	return filepath.Join(filepath.Dir(e.rootFile.internalPath), imgRelativePath), nil
}

// Attempts to convert provided image data to required format before
// setting epub field
func (e *Epub) SetCoverImage(cover []byte) error {

	newMediaType := http.DetectContentType(cover)
	if !strings.HasPrefix(newMediaType, "image/") {
		return fmt.Errorf("invalid media type %s", newMediaType)
	}

	destination, err := e.getCoverPath()
	if err != nil {
		return err
	}

	reqMediaType := filepath.Ext(destination)
	img, _, err := image.Decode(bytes.NewReader(cover))
	if err != nil {
		return err
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = nil
	switch reqMediaType {
	case ".jpg":
	case ".jpeg":
		err = jpeg.Encode(writer, img, nil)
	case ".png":
		err = png.Encode(writer, img)
	case ".gif":
		err = gif.Encode(writer, img, nil)
	default:
		err = fmt.Errorf("image conversion to %s is not supported", reqMediaType)
	}

	writer.Flush()

	if err != nil {
		return err
	}

	e.coverImage = b.Bytes()
	return nil
}
