package book

import (
	"nubayrah/epub"

	"github.com/google/uuid"
)

type Book struct {
	ID uuid.UUID `json:"id"`
	epub.Metadata
	Filepath string `json:"filePath"`
}

type Books []*Book
