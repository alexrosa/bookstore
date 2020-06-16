package books

import (
	"encoding/json"
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

func (handler *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := handler.repository.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(books)
}
