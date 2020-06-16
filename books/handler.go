package books

import (
	"net/http"
)

type BookHandler struct {
	repository BooksRepository
}

func NewBookHandler(repository BooksRepository) *BookHandler {
	return &BookHandler{
		repository: repository,
	}
}

func (bh *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

}
