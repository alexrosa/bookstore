package server

import (
	"database/sql"
	"net/http"

	"github.com/bookstore/db"

	"github.com/bookstore/books"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var dbConn *sql.DB

func init() {
	dbConn = db.GetDBConnection()
}

func StartService() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	//book router
	bookRouter(r)
	return r
}

func bookRouter(r chi.Router) {
	repository := books.NewBooksRepository(dbConn)
	bookHandler := books.NewBookHandler(*repository)
	r.Get("/books", bookHandler.GetBooks)

}
