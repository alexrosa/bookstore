package books

import (
	"context"
	"database/sql"
)

//go:generate mockery -name=BooksStorage .
type BooksStorage interface {
	GetByID(ctx context.Context, id int64) (*Book, error)
	ListAll(ctx context.Context) ([]Book, error)
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book Book) error
	Delete(ctx context.Context, id int64) error
}

type BooksRepository struct {
	dbConn *sql.DB
}

func NewBooksRepository(conn *sql.DB) *BooksRepository {
	return &BooksRepository{dbConn: conn}
}

func (br *BooksRepository) GetByID(ctx context.Context, id int64) (Book, error) {
	var book Book

	defer br.dbConn.Close()

	err := br.dbConn.QueryRowContext(ctx, `SELECT id, title, author, issued, pages, last_updated, created_at FROM books WHERE id =?`, id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Issued, &book.Pages, &book.LastUpdated, &book.CreatedAt)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (br *BooksRepository) ListAll(ctx context.Context) ([]Book, error) {
	selectStmt := `SELECT id, title, author, issued, pages, last_updated, created_at FROM books;`
	var books []Book

	rows, err := br.dbConn.QueryContext(ctx, selectStmt)
	if err != nil {
		return books, nil
	}
	defer rows.Close()

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Issued, &b.Pages, &b.LastUpdated, &b.CreatedAt); err != nil {
			return books, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (br *BooksRepository) Create(ctx context.Context, book *Book) (err error) {
	insertStmt := `INSERT INTO books (title, author, issued, pages) VALUES (?, ?, ?, ?);`

	var tx *sql.Tx
	tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, insertStmt, book.Title, book.Author, book.Issued, book.Pages)
	if err != nil {
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return
	}
	book.ID = lastID
	return
}

func (br *BooksRepository) Update(ctx context.Context, book Book) (err error) {
	updateStmt := `UPDATE books SET title=?, author=?, issued=?, pages=? WHERE id=?;`

	var tx *sql.Tx
	tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, updateStmt, book.Title, book.Author, book.Issued, book.Pages, book.ID)

	if err != nil {
		return
	}
	return
}

func (br *BooksRepository) Delete(ctx context.Context, id int64) (err error) {
	deleteStmt := `DELETE FROM books WHERE id=?;`

	var tx *sql.Tx
	tx, err = br.dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		return
	}

	return
}
