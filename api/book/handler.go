package book

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nubayrah/epub"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	defer file.Close()

	epubObj, err := epub.Import(file)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("error opening epub archive %v", err)
	}

	book, err := a.repository.Create(&Book{
		Metadata: *epubObj.ExtractMetadata(),
		ID:       uuid.New(),
		Filepath: epubObj.FilePath,
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

// Handler for getting a specific book at /books/{bookID}
func (a *API) handleGetBookCover(w http.ResponseWriter, r *http.Request) {

	// Grab UUID from url
	id := chi.URLParam(r, "id")
	UUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error parsing uuid from url: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Read item from database
	book, err := a.repository.Read(UUID)
	if err != nil {
		log.Printf("error finding book in db: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the filePath for the item
	e, err := epub.OpenEpub(book.Filepath)
	if err != nil {
		log.Printf("error parsing path for item %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get CoverFile file object
	file, err := e.GetCoverFile()
	if err != nil {
		log.Printf("error creating file object for item %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Open a io.Reader for the object
	fileReader, err := file.Open()
	if err != nil {
		log.Printf("error creating reader for object %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Copies content of the file to the response writer
	_, err = io.Copy(w, fileReader)
	if err != nil {
		log.Printf("error copying content to response %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Handler for Deleting a specific book.
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
