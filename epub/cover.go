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
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Returns internal path to cover image
func (e *Epub) GetCoverPath() (string, error) {
	coverId := e.RootFile.getCoverId()
	if coverId == "" {
		return "", fmt.Errorf("cover id not found in root doc")
	}

	itemElem := e.RootFile.FindElement(fmt.Sprintf("//package/manifest/item[@id='%s']", coverId))
	if itemElem == nil {
		return "", errors.New("cover image item not found in manifest")
	}

	imgRelativePath := itemElem.SelectAttrValue("href", "")

	return filepath.Join(filepath.Dir(e.RootFile.internalPath), imgRelativePath), nil
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

// Gets and returns *zip.File pointing to the coverFile path.
func (e *Epub) GetCoverFile() (*zip.File, error) {
	path, err := e.GetCoverPath()
	if err != nil {
		return nil, err
	}
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

	return file, nil
}
