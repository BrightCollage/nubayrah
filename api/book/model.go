package book

import "nubayrah/epub"

type Book struct {
	epub.Metadata
	ID       string `json:"id"`
	Filepath string
}
