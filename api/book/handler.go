package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nubayrah/epub"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	config "github.com/spf13/viper"
	"gorm.io/gorm"
)

type API struct {
	repository *Repository
}

func NewAPI(db *gorm.DB) *API {
	return &API{
		repository: NewRepository(db),
	}
}

// Handler for importing an epub
func (a *API) handleImportBook(w http.ResponseWriter, r *http.Request) {
	// Read file contents from request
	file, _, err := r.FormFile("epub")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error reading file from post requiest %v", err)
		return
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error copying file data %v", err)
		return
	}

	data := buf.Bytes()

	// Validate first by checking magic bytes, then attempting to parse the epub's metadata
	var magic [4]byte
	copy(magic[:], data[:4])
	if !epub.CheckMagic(magic) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("magic bytes invalid for uploaded epub: %v", magic)
		return
	}

	epubObj, err := epub.OpenEpubBytes(data)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("error opening epub archive %v", err)
	}

	// Write epub to disk at library/author/title.epub
	targetDir := filepath.Join(config.GetString("library_path"), epubObj.Metadata.Author)
	targetDir = sanitizeDirName(targetDir)

	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot create directories %v", err)
		return
	}

	targetFile := filepath.Join(targetDir, sanitizeFileName(epubObj.Metadata.Title)) + ".epub"
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
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("unable to find unused filename for %v", targetFile)
				return
			}
		}
	}

	epubObj.Filepath = targetFile
	err = os.WriteFile(targetFile, data, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error writing epub file to disk %v", err)
	}

	book, err := a.repository.Create(&Book{
		Metadata: *epubObj.ExtractMetadata(),
		ID:       uuid.New(),
		Filepath: targetDir,
	})
	if err != nil {
		log.Printf("error writing books into database %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(book)
	if err != nil {
		log.Printf("error marshalling books into json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Extract cover image to library/author/{covername}.ext
	// coverFilepath, err := epub.ExtractCoverImage(targetDir)
	// if err != nil {
	// 	log.Printf("error extracting cover image %v", err)
	// }

	// // Rename cover image to n.ext where n is the id for the new row in the db
	// n, _ := res.LastInsertId()
	// d, f := filepath.Split(coverFilepath)
	// ext := filepath.Ext(f)
	// err = os.Rename(coverFilepath, filepath.Join(d, strconv.FormatInt(n, 10)+ext))
	// if err != nil {
	// 	log.Printf("error renaming cover image %v", err)
	// }

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// Handler for root link /books
func (a *API) handleGetAllBooks(w http.ResponseWriter, _ *http.Request) {

	books, err := a.repository.List()

	if err != nil {
		log.Printf("error reading rows %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(books)
	if err != nil {
		log.Printf("error marshalling books into json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

// Handler for getting a specific book at /books/{bookID}
func (a *API) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error parsing uuid from url: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	book, err := a.repository.Read(UUID)
	if err != nil {
		log.Printf("error finding book in db: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j, err := json.Marshal(book)
	if err != nil {
		log.Printf("error marshalling books into json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

// Handler for getting a specific book.
func (a *API) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	// Grab ID from the URL, which is /todo/{todoID}
	id := chi.URLParam(r, "id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error when parsing UUID from url: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	count, err := a.repository.Delete(UUID)
	if err != nil {
		log.Printf("error when deleting UUID from DB: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Count deleted: %v", count)

	w.WriteHeader(http.StatusNoContent)
}

// // Scan and parse database row into book
// // Some fields (eg any array) cannot be stored directly in the db and are
// // encoded/decoded as json
// func rowToMetadata(row *sql.Rows) (*Book, error) {
// 	b := &Book{}

// 	var subjects string
// 	var contribs string

// 	err := row.Scan(
// 		&b.ID,
// 		&b.Filepath,
// 		&b.Metadata.Title,
// 		&b.Metadata.TitleSort,
// 		&b.Metadata.Author,
// 		&b.Metadata.AuthorSort,
// 		&b.Metadata.Language,
// 		&b.Metadata.Series,
// 		&b.Metadata.SeriesNum,
// 		&subjects,
// 		&b.Metadata.Isbn,
// 		&b.Metadata.Publisher,
// 		&b.Metadata.PubDate,
// 		&b.Metadata.Rights,
// 		&contribs,
// 		&b.Metadata.Description,
// 		&b.Metadata.Uid,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	json.Unmarshal([]byte(subjects), &b.Subjects)
// 	json.Unmarshal([]byte(contribs), &b.Contributors)

// 	return b, nil
// }
