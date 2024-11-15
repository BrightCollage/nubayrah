// The data Models and conversions as needed.

package book

import (
	"nubayrah/epub"

	"github.com/google/uuid"
)

type Book struct {
	ID uuid.UUID `json:"id" gorm:"<-:create"`
	epub.Metadata
	Filepath string `json:"filePath"`
}

type Books []*Book
