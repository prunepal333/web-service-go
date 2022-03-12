package main

import "math/rand"
import "github.com/gin-gonic/gin"
import "net/http"
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/go-sql-driver/mysql"
)
//➔ books: structure that holds information about book records.
type Book struct {
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Page   int32   `json:"page"`
}

var db *sql.DB

var books []Book

func generate_rand_digit_str() string {
	return fmt.Sprintf("249804%d48493", rand.Intn(100))
}
func getBooks(c *gin.Context) {
	rows, err := db.Query("SElECT * FROM books")
	if err != nil{
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}
	defer rows.Close()

	//loop through rows using Scan
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ISBN, &book.Title, &book.Author, &book.Price, &book.Page); err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return			
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H {"message": "Internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func postBooks(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(400, gin.H{"error": "Bad input!"})
		return
	}
	escapedQuery := "INSERT INTO books (isbn, title, author, price, page) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(escapedQuery, newBook.ISBN, newBook.Title, newBook.Author, newBook.Price, newBook.Page)
	
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	//➔ Add the new book to the slice
	books = append(books, newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

func notFoundError(c *gin.Context){
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Not Found!",
	})
}

func getBookByID(c *gin.Context){

	id := c.Param("id")

	row := db.QueryRow("SELECT * FROM books WHERE isbn=?", id)

	if err := row.Scan(&book.ISBN, &book.Title, &book.Author, &book.Price, &book.Page) ; err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, {"message": "Not found"})
			return
		}
		c.JSON(500, {"message": "Internal server error"})
	}
	
	for _, book := range books {
		if book.ISBN == id {
			c.IndentedJSON(http.StatusOK, book)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}


func updateBookByID(c *gin.Context) {
	var updatedBook Book
	id := c.Param("id")
	for i := 0; i < len(books); i++ {
		book := &books[i]
		if book.ISBN == id {
			if err := c.ShouldBindJSON(&updatedBook); err != nil {
				c.JSON(400, gin.H{"error": "Bad input!"})
				return
			}
			book.ISBN = updatedBook.ISBN
			book.Title = updatedBook.Title
			book.Author = updatedBook.Author
			book.Price = updatedBook.Price
			book.Page = updatedBook.Page
			c.JSON(200, gin.H{"message": "book info updated successfully!"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	
}

func main() {

	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "db-rest-demo",
	}
	//➔ Explicit declaration of variable
	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	router := gin.Default()

	router.GET("/", notFoundError)
	router.GET("/books", getBooks)
	router.POST("/books", postBooks)

	router.GET("/books/:id", getBookByID)

	router.PUT("/books/:id", updateBookByID)
	router.Run("localhost:8080")
}