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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/go-xmlfmt/xmlfmt"
	"golang.org/x/net/html/charset"
)

type Epub struct {
	Filepath   string
	Metadata   *Metadata
	fileHandle *zip.Reader
	rootFile   *RootFile
	coverImage []byte
}

// Opens and parses epub from file on disk
func OpenEpub(filepath string) (*Epub, error) {
	epub := &Epub{}
	epub.Filepath = filepath

	err := epub.Reload()
	if err != nil {
		return nil, err
	}

	return epub, nil
}

// Opens and parses epub archive from byte array
func OpenEpubBytes(data []byte) (*Epub, error) {
	var err error = nil
	epub := &Epub{}

	rdr := bytes.NewReader(data)

	epub.fileHandle, err = zip.NewReader(rdr, rdr.Size())
	if err != nil {
		return nil, err
	}

	err = epub.Reload()
	if err != nil {
		return nil, err
	}

	return epub, nil
}

// Checks if first 4 bytes match epub magic bytes described here:
// https://en.wikipedia.org/wiki/List_of_file_signatures
// Does not confirm that a file *is* an epub or valid archive
func CheckMagic(data [4]byte) bool {
	return data == [4]byte{0x50, 0x4B, 0x03, 0x04} ||
		data == [4]byte{0x50, 0x4B, 0x05, 0x06} ||
		data == [4]byte{0x50, 0x4B, 0x07, 0x08}
}

// Reopens and reads epub from disk
func (e *Epub) Reload() error {
	var err error = nil

	if e.Filepath == "" {
		log.Print("epub filepath is empty, cannot reload from disk")
	} else {
		e.Close()

		bts, err := os.ReadFile(e.Filepath)
		if err != nil {
			return err
		}
		rdr := bytes.NewReader(bts)

		e.fileHandle, err = zip.NewReader(rdr, rdr.Size())
		if err != nil {
			return err
		}
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

	data, err := io.ReadAll(c)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Extract cover image into destination directory
// eg ExtractCoverImage("/library/author") may result in a file
// `/library/author/cover_image_123456789.png`
// Returns the full path to the cover image
func (e *Epub) ExtractCoverImage(destDir string) (string, error) {
	source, err := e.GetCoverPath()
	if err != nil {
		return "", err
	}

	data, err := e.ReadFile(source)
	if err != nil {
		return "", err
	}

	destFile := filepath.Join(destDir, filepath.Base(source))

	err = os.WriteFile(destFile, data, os.ModePerm)
	if err != nil {
		return "", nil
	}

	return destFile, nil
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
	tmpDir := filepath.Join(os.TempDir(), filepath.Base(strings.Split(e.Filepath, ".")[0]))
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
		coverPath, err := e.GetCoverPath()
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
	e.fileHandle = nil
	err = os.Rename(tmpFile, e.Filepath)
	if err != nil {
		return err
	}

	return nil
}

// Returns internal path to cover image
func (e *Epub) GetCoverPath() (string, error) {
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

	destination, err := e.GetCoverPath()
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
