package main

// Первая лаба на GO Реализовать API для библиотеки:
// реализовать ручку для получения книг и дату их возврата по читателю (книги которые на руках);
// реализовать ручку для выдачи книги на чтение с ограниченным сроком (срок выдачи книги на 2 недели);
// ready - реализовать ручку добавления читателя (имя/номер/адрес/дата рождения);
// ready - реализовать ручку добавления книги (Название/кол-во/автор/год выпуска).

// Одновременно читатель не может взять больше 3 книг.
// Первая ручка должна возвращать отдельно просроченные или не просроченные книги.
// Базы данных нет - данные хранятся в переменных.

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var books = []book{
	{"Война и мир", 5, "Лев Толстой", newYear(1869)},
	{"Преступление и наказание", 3, "Федор Достоевский", newYear(1866)},
	{"Мастер и Маргарита", 4, "Михаил Булгаков", newYear(1967)},
	{"Анна Каренина", 2, "Лев Толстой", newYear(1877)},
	{"Евгений Онегин", 6, "Александр Пушкин", newYear(1833)},
	{"Идиот", 3, "Федор Достоевский", newYear(1869)},
	{"Тихий Дон", 4, "Михаил Шолохов", newYear(1940)},
	{"Доктор Живаго", 2, "Борис Пастернак", newYear(1957)},
	{"Отцы и дети", 5, "Иван Тургенев", newYear(1862)},
	{"Мертвые души", 3, "Николай Гоголь", newYear(1842)},
}

var readers = []reader{
	{
		Name:        "Мария Петрова",
		PhoneNumber: "002",
		Adress:      "ул. Советская, д.3, кв.12",
		DateOfBirth: time.Date(1990, time.July, 24, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
	{
		Name:        "Ольга Кузнецова",
		PhoneNumber: "004",
		Adress:      "ул. Пушкина, д.5, кв.2",
		DateOfBirth: time.Date(1982, time.November, 10, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
	{
		Name:        "Дмитрий Орлов",
		PhoneNumber: "005",
		Adress:      "ул. Гагарина, д.18, кв.7",
		DateOfBirth: time.Date(1995, time.May, 15, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
	{
		Name:        "Николай Федоров",
		PhoneNumber: "007",
		Adress:      "ул. Толстого, д.4, кв.9",
		DateOfBirth: time.Date(1972, time.December, 22, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
	{
		Name:        "Павел Новиков",
		PhoneNumber: "009",
		Adress:      "ул. Горького, д.7, кв.3",
		DateOfBirth: time.Date(1980, time.August, 1, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
	{
		Name:        "Анна Сергеева",
		PhoneNumber: "010",
		Adress:      "пр. Карла Маркса, д.16, кв.6",
		DateOfBirth: time.Date(1992, time.April, 29, 0, 0, 0, 0, time.UTC),
		bookInUse:   []bookInUse{},
	},
}

type book struct {
	Title  string
	Copies int
	Author string
	Issue  time.Time
}

// postBooks adds an album from JSON received in the request body.
func postBooks(c *gin.Context) {
	var newBook book

	// Call BindJSON to bind the received JSON to newBook.
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	// Add the new book to the slice.
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// getBooks responds with the list of all books as JSON.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// getBookByTitle locates the album whose Title value matches the title
// parameter sent by the client, then returns that book as a response.
func getBookByTitle(c *gin.Context) {
	title := c.Param("Title")

	// Loop over the list of books, looking for
	// an book whose Title value matches the parameter.
	for _, a := range books {
		if a.Title == title {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

type reader struct {
	Name        string
	PhoneNumber string
	Adress      string
	DateOfBirth time.Time
	bookInUse   []bookInUse
}

type bookInUse struct {
	NameOfBook book
	DateOfRent time.Time
}

// func getBookById(c *gin.Context, id int)  {
// 	c.IndentedJSON()
// }

// getBooks responds with the list of all books as JSON.
func getReaders(c *gin.Context) {
	c.JSON(http.StatusOK, readers)
}

func newYear(year int) time.Time {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/book/:Title", getBookByTitle)
	router.GET("/readers", getReaders)
	router.POST("/books", postBooks)
	// router.GET('/readers/:Name', getReaderBooks)
	// router.PATCH("/readers/", rentBookByTitle)
	// router.PATCH("/readers/:Name", returnBookByTitle)

	router.Run("localhost:8080")
}
