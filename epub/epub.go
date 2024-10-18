package epub

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/go-xmlfmt/xmlfmt"
	config "github.com/spf13/viper"
	"golang.org/x/net/html/charset"
)

type Epub struct {
	FilePath   string // Path to .epub
	FileDir    string // Path to directory holding .epub
	FileName   string
	Metadata   *Metadata
	fileHandle *zip.Reader
	RootFile   *RootFile
	coverImage []byte
}

// Opens and parses epub from file on disk
func OpenEpub(path string) (*Epub, error) {
	epub := &Epub{}
	epub.FilePath = path
	epub.FileDir = filepath.Dir(path)

	err := epub.Reload()
	if err != nil {
		return nil, err
	}

	return epub, nil
}

// Opens and parses epub from multi-part file from request and then saves it to disk
// at config.library_path/author/title.epub
func Import(file multipart.File) (*Epub, error) {

	// Create bytes buffer from file.
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()

	// Check magic bytes to ensure epub.
	err = checkMagic(data)
	if err != nil {
		return nil, err
	}

	e := &Epub{}

	rdr := bytes.NewReader(data)

	e.fileHandle, err = zip.NewReader(rdr, rdr.Size())
	if err != nil {
		return nil, err
	}

	err = e.Load()
	if err != nil {
		return nil, err
	}

	// Write epub to disk at library/author/title.epub
	targetDir := filepath.Join(config.GetString("library_path"), e.Metadata.Author)
	targetDir = sanitizeDirName(targetDir)

	e.FileDir = targetDir

	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		log.Printf("cannot create directories %v", err)
		return nil, err
	}

	targetFile := filepath.Join(targetDir, sanitizeFileName(e.Metadata.Title)) + ".epub"
	// If a file exists with the desired name, start incrementing as filename_1
	// until an unused filename is found
	if fileExists(targetFile) {
		k := strings.LastIndex(targetFile, ".")
		targetFile = targetFile[:k] + "_%d" + ".epub"
		for i := 1; i < 256; i++ {
			numberedTarget := fmt.Sprintf(targetFile, i)
			if !fileExists(numberedTarget) {
				targetFile = numberedTarget
				break
			}
			if i == 255 {
				return nil, errors.New("unable to find unused filename")
			}
		}
	}

	e.FilePath = targetFile
	err = os.WriteFile(targetFile, data, os.ModePerm)
	if err != nil {
		return nil, errors.New("unable to write to disk")
	}
	// Set the FileName
	e.FileName = filepath.Base(e.FilePath)

	return e, nil
}

// Checks if first 4 bytes match epub magic bytes described here:
// https://en.wikipedia.org/wiki/List_of_file_signatures
// Does not confirm that a file *is* an epub or valid archive
func checkMagicBytes(data [4]byte) bool {
	return data == [4]byte{0x50, 0x4B, 0x03, 0x04} ||
		data == [4]byte{0x50, 0x4B, 0x05, 0x06} ||
		data == [4]byte{0x50, 0x4B, 0x07, 0x08}
}

func checkMagic(data []byte) error {

	// Validate first by checking magic bytes, then attempting to parse the epub's metadata
	var magic [4]byte
	copy(magic[:], data[:4])

	if !checkMagicBytes(magic) {
		return errors.New("magic byte not found")
	}

	return nil
}

// Reopens and reads epub from disk
func (e *Epub) Reload() error {

	if e.FilePath == "" {
		log.Print("epub filepath is empty, cannot reload from disk")
	} else {
		e.Close()

		bts, err := os.ReadFile(e.FilePath)
		if err != nil {
			return err
		}
		rdr := bytes.NewReader(bts)

		e.fileHandle, err = zip.NewReader(rdr, rdr.Size())
		if err != nil {
			return err
		}
	}

	return e.Load()
}

// Loads epub into Epub object.
func (e *Epub) Load() error {
	var err error = nil

	e.RootFile, err = e.getRootFile()
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
	e.Metadata = e.ExtractMetadata()
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
	tmpDir := filepath.Join(os.TempDir(), filepath.Base(strings.Split(e.FilePath, ".")[0]))
	tmpFile := tmpDir + ".epub"

	e.unpack(tmpDir)

	defer os.Remove(tmpDir)

	// RootFile + Metadata
	e.RootFile.InsertMetadata(e.Metadata)
	rfStr, err := e.RootFile.WriteToString()
	if err != nil {
		return err
	}
	rfStr = xmlfmt.FormatXML(rfStr, "", "  ")

	err = os.WriteFile(filepath.Join(tmpDir, e.RootFile.internalPath), []byte(rfStr), os.ModePerm)
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
	err = os.Rename(tmpFile, e.FilePath)
	if err != nil {
		return err
	}

	return nil
}
