package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/epub"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Book struct {
	epub.Metadata
	ID       int `json:"id"`
	Filepath string
}

// Handler for importing an epub
func handleImportBook(w http.ResponseWriter, r *http.Request) {
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

	epub, err := epub.OpenEpubBytes(data)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("error opening epub archive %v", err)
	}

	// Write epub to disk at library/author/title.epub
	targetDir := filepath.Join(userConfig.LibraryRoot, epub.Metadata.Author)
	targetDir = sanitizeDirName(targetDir)

	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot create directories %v", err)
		return
	}

	targetFile := filepath.Join(targetDir, sanitizeFileName(epub.Metadata.Title)) + ".epub"
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

	epub.Filepath = targetFile
	err = os.WriteFile(targetFile, data, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error writing epub file to disk %v", err)
	}

	// Extract cover image to library/author/{covername}.ext
	coverFilepath, err := epub.ExtractCoverImage(targetDir)
	if err != nil {
		log.Printf("error extracting cover image %v", err)
	}

	// Insert into db
	subjects, err := json.Marshal(epub.Metadata.Subjects)
	if err != nil {
		subjects = []byte{}
	}

	contribs, err := json.Marshal(epub.Metadata.Contributors)
	if err != nil {
		contribs = []byte{}
	}

	insertQ := `INSERT INTO library
	(id,  filepath, title, titleSort, author, authorSort, language, series, seriesNum, subjects, isbn, publisher, pubDate, rights, contributors, description, uid) VALUES
	(NULL,       $1,   $2,        $3,     $4,         $5,       $6,     $7,        $8,       $9,  $10,       $11,     $12,    $13,          $14,         $15, $16)`
	res, err := DB.Exec(insertQ,
		epub.Filepath,
		epub.Metadata.Title,
		epub.Metadata.TitleSort,
		epub.Metadata.Author,
		epub.Metadata.AuthorSort,
		epub.Metadata.Language,
		epub.Metadata.Series,
		strconv.FormatFloat(epub.Metadata.SeriesNum, 'f', 2, 64),
		string(subjects),
		epub.Metadata.Isbn,
		epub.Metadata.Publisher,
		epub.Metadata.PubDate,
		epub.Metadata.Rights,
		string(contribs),
		epub.Metadata.Description,
		epub.Metadata.Uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error when querying to DB %v", err)
		os.Remove(coverFilepath)
		os.Remove(epub.Filepath)
		return
	}

	// Rename cover image to n.ext where n is the id for the new row in the db
	n, _ := res.LastInsertId()
	d, f := filepath.Split(coverFilepath)
	ext := filepath.Ext(f)
	err = os.Rename(coverFilepath, filepath.Join(d, strconv.FormatInt(n, 10)+ext))
	if err != nil {
		log.Printf("error renaming cover image %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler for root link /books
func handleGetAllBooks(w http.ResponseWriter, _ *http.Request) {
	rows, err := DB.Query("SELECT * from library;")
	if err != nil {
		log.Printf("error reading database %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	books := []Book{}
	for rows.Next() {
		b, err := rowToMetadata(rows)
		if err != nil {
			log.Printf("error reading database %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		books = append(books, *b)
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
func handleGetBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error parsing %d from string to integer %v", bookID, err)
		return
	}

	row, err := DB.Query("SELECT * FROM library WHERE id=$1;", bookID)
	if err != nil {
		log.Printf("error reading database %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	row.Next()
	book, err := rowToMetadata(row)
	if err != nil {
		log.Printf("error reading database %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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
func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	// Grab ID from the URL, which is /todo/{todoID}
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error parsing %d from string to integer %v", bookID, err)
		return
	}

	// SQL query to delete a row.
	sqlStatement := `
	DELETE FROM library
	WHERE id = $1;`
	res, err := DB.Exec(sqlStatement, bookID)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Printf("Count deleted: %v", count)

	w.WriteHeader(http.StatusNoContent)
}

// Scan and parse database row into book
// Some fields (eg any array) cannot be stored directly in the db and are
// encoded/decoded as json
func rowToMetadata(row *sql.Rows) (*Book, error) {
	b := &Book{}

	var subjects string
	var contribs string

	err := row.Scan(
		&b.ID,
		&b.Filepath,
		&b.Metadata.Title,
		&b.Metadata.TitleSort,
		&b.Metadata.Author,
		&b.Metadata.AuthorSort,
		&b.Metadata.Language,
		&b.Metadata.Series,
		&b.Metadata.SeriesNum,
		&subjects,
		&b.Metadata.Isbn,
		&b.Metadata.Publisher,
		&b.Metadata.PubDate,
		&b.Metadata.Rights,
		&contribs,
		&b.Metadata.Description,
		&b.Metadata.Uid,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(subjects), &b.Subjects)
	json.Unmarshal([]byte(contribs), &b.Contributors)

	return b, nil
}
