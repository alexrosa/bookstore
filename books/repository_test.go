package books_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/bookstore/books"

	"github.com/stretchr/testify/assert"
)

const (
	InsertQuery = `INSERT INTO books `
	UpdateQuery = `UPDATE books `
	DeleteQuery = `DELETE FROM books`
	SelectQuery = `SELECT (.+) FROM books;`
)

func TestBooksRepository(t *testing.T) {
	ctx := context.Background()
	newBook := books.Book{
		Title:  "Go Development",
		Author: "Alex Rosa",
		Issued: 2010,
		Pages:  500,
	}

	t.Run("save", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.Nil(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec(InsertQuery).WithArgs(newBook.Title, newBook.Author, newBook.Issued, newBook.Pages).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
		mock.ExpectCommit()

		repo := books.NewBooksRepository(db)
		err = repo.Create(ctx, &newBook)

		assert.Nil(t, err)
		assert.Equal(t, newBook.ID, int64(1))

		//check if everything worked well
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("save-error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec(InsertQuery).WithArgs(newBook.Title, newBook.Author, newBook.Issued, newBook.Pages).WillReturnError(fmt.Errorf("bang"))
		mock.ExpectRollback()

		repo := books.NewBooksRepository(db)
		err = repo.Create(ctx, &newBook)

		assert.Error(t, err)
		assert.Equal(t, "bang", err.Error())

		//check if everything worked well
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("update", func(t *testing.T) {
		updatedBook := books.Book{
			ID:          1,
			Title:       "some_title",
			Author:      "some_author",
			Issued:      1999,
			Pages:       1,
			LastUpdated: time.Now(),
			CreatedAt:   time.Now(),
		}

		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec(UpdateQuery).WithArgs(updatedBook.Title, updatedBook.Author,
			updatedBook.Issued, updatedBook.Pages, updatedBook.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := books.NewBooksRepository(db)

		err = repo.Update(ctx, updatedBook)
		assert.Nil(t, err)

		//checking if everything worked well
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()
		mock.ExpectBegin()
		mock.ExpectExec(DeleteQuery).WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := books.NewBooksRepository(db)

		err = repo.Delete(ctx, 1)
		assert.Nil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("listBooks", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()
		now := time.Now()
		expectedList := []books.Book{
			{
				ID:          1,
				Title:       "test",
				Author:      "Alex",
				Issued:      2010,
				Pages:       10,
				LastUpdated: now,
				CreatedAt:   now,
			},
		}
		rows := sqlmock.NewRows([]string{"id", "title", "author", "issued", "pages", "last_updated", "created_at"}).
			AddRow(1, "test", "Alex", 2010, 10, now, now)
		mock.ExpectQuery(SelectQuery).WillReturnRows(rows)
		repo := books.NewBooksRepository(db)
		list, err := repo.ListAll(ctx)

		assert.Nil(t, err)
		assert.Equal(t, expectedList, list)
	})

	t.Run("GetByID", func(t *testing.T) {
		selectQuery := `SELECT id, title, author, issued, pages, last_updated, created_at FROM books WHERE id =?`
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()
		now := time.Now()
		expectedBook := books.Book{
			ID:          1,
			Title:       "test",
			Author:      "Alex",
			Issued:      2010,
			Pages:       10,
			LastUpdated: now,
			CreatedAt:   now,
		}
		rows := sqlmock.NewRows([]string{"id", "title", "author", "issued", "pages", "last_updated", "created_at"}).
			AddRow(1, "test", "Alex", 2010, 10, now, now)
		mock.ExpectQuery(selectQuery).WithArgs(int64(1)).WillReturnRows(rows)

		repo := books.NewBooksRepository(db)
		resultBook, err := repo.GetByID(ctx, int64(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedBook, resultBook)
	})
}
