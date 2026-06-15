package products

import (
	"context"
	// "log"
)

type Service interface {
	GetAllProducts(ctx context.Context) (products []Book, err error)
}

type svc struct {
	// database
}

func NewService() Service {
	return &svc{}
}
	type Book struct {
		book_id     int
		title       string
		author      string
		description string
	}
func (svc *svc) GetAllProducts(ctx context.Context) (products []Book, err error) {
	// Database Logic
	// Get Books from Database



	// query := `SELECT book_id, title, author, description FROM books`

	// rows, err := db.Query(query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var books []Book
	// for rows.Next() {
	// 	var book Book
	// 	rows.Scan(&book.book_id, &book.title, &book.author, &book.description)
	// 	books = append(books, book)
	// }

	// if err := rows.Error(); err != nil {
	// 	log.Fatal(err)
	// }

	books := []Book{
		{1, "", "", ""},
	}


	return books, nil

}
