package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ImportBookBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Handler for creating importing a book.
func handleImportBook(w http.ResponseWriter, r *http.Request) {

	var body ImportBookBody

	// Creates a Decoder which reads json data from a data stream and converts it to structured data.
	// NewDecoder takes in an io.Reader, r.Body in this case, which then reads the ioStream.
	// .Decodes(&dst) takes an argument, which is the destination location for the decoded output information.
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body into ImportBookBody struct %v", err)
		return
	}

	// We ask the DB object to send a SQL query to INSERT a data row to the Books table with 2 columns: title and description.
	// The values are the $1 -> body.Title and $2 -> body.Description
	if err := DB.QueryRow("INSERT INTO library (title, description) VALUES ($1, $2)", body.Title, body.Description).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error when querying to DB %v", err)
		return
	}

	// Once done, we just return a status saying that the create was successful.
	w.WriteHeader(http.StatusCreated)

}

// Handler for root link /books
func handleGetAllBooks(w http.ResponseWriter, _ *http.Request) {
	var Books []Book

	// Query the database table 'Books' and select the columns: id, title, and description from the table
	rows, err := DB.Query("SELECT id, title, description FROM library")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// For loop checks the next item in row, until there are no rows left, in which it will return a False.
	for rows.Next() {
		// fetchedItem is a struct object that contains json marshalling format for our Book.
		var fetchedItem Book
		// rows.Scan(...) will fill each argument with a value it finds in the row.
		// That is why we pass them the location (pointer) to where we want the values filled. In this case, it is
		// each of the struct's members.
		if err := rows.Scan(&fetchedItem.ID, &fetchedItem.Title, &fetchedItem.Description); err != nil {
			// Check any line for error.
			log.Printf("Error when scanning rows %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// If columns are successfully 'Scanned' then we add the object to the list.
		Books = append(Books, fetchedItem)
	}

	// Takes the list of object we have above and we will make it json data.
	j, err := json.Marshal(Books)
	// Check error in marshalling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling Books into json %v", err)
		return
	}

	// If query from database + creating marshallable object + converting to json object OK, then we write the values back.
	w.Write(j)
}

// Handler for getting a specific book.
func handleGetBook(w http.ResponseWriter, r *http.Request) {
	// Grab ID from the URL, which is /todo/{todoID}
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error parsing %d from string to integer %v", bookID, err)
		return
	}

	var Book Book

	// Query to postgresql that grabs the columns, id, title, description from the table called Books and filters columns where the id is
	// input value.
	query := `SELECT id, title, description FROM library WHERE id=$1;`
	// Query for a single row
	row := DB.QueryRow(query, bookID)
	if err := row.Scan(&Book.ID, &Book.Title, &Book.Description); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Error when scanning rows: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	j, err := json.Marshal(Book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while marshalling Book to json %v", err)
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
