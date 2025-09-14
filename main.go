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

const timeToExpire = 14 * 24 * time.Hour

var books = []book{
	{"Война и мир", 5, "Лев Толстой", newYear(1869)},
	{"Преступление и наказание", 3, "Федор Достоевский", newYear(1866)},
	{"Мастер и Маргарита", 4, "Михаил Булгаков", newYear(1967)},
	{"Анна Каренина", 2, "Лев Толстой", newYear(1877)},
	{"Евгений Онегин", 6, "Александр Пушкин", newYear(1833)},
	{"Идиот", 0, "Федор Достоевский", newYear(1869)},
	{"Тихий Дон", 1, "Михаил Шолохов", newYear(1940)},
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
		booksInUse: []booksInUse{
			{
				NameOfBook: books[1],                      // например, "Преступление и наказание"
				DateOfRent: time.Now().AddDate(0, 0, -20), // Взята 20 дней назад (> 2 недель)
			},
			{
				NameOfBook: books[2],                     // например, "Мастер и Маргарита"
				DateOfRent: time.Now().AddDate(0, 0, -5), // Взята 5 дней назад (< 2 недель)
			},
		},
	},
	{
		Name:        "Ольга Кузнецова",
		PhoneNumber: "004",
		Adress:      "ул. Пушкина, д.5, кв.2",
		DateOfBirth: time.Date(1982, time.November, 10, 0, 0, 0, 0, time.UTC),
		booksInUse:  []booksInUse{},
	},
	{
		Name:        "Дмитрий Орлов",
		PhoneNumber: "005",
		Adress:      "ул. Гагарина, д.18, кв.7",
		DateOfBirth: time.Date(1995, time.May, 15, 0, 0, 0, 0, time.UTC),
		booksInUse:  []booksInUse{},
	},
	{
		Name:        "Николай Федоров",
		PhoneNumber: "007",
		Adress:      "ул. Толстого, д.4, кв.9",
		DateOfBirth: time.Date(1972, time.December, 22, 0, 0, 0, 0, time.UTC),
		booksInUse:  []booksInUse{},
	},
	{
		Name:        "Павел Новиков",
		PhoneNumber: "009",
		Adress:      "ул. Горького, д.7, кв.3",
		DateOfBirth: time.Date(1980, time.August, 1, 0, 0, 0, 0, time.UTC),
		booksInUse:  []booksInUse{},
	},
	{
		Name:        "Анна Сергеева",
		PhoneNumber: "010",
		Adress:      "пр. Карла Маркса, д.16, кв.6",
		DateOfBirth: time.Date(1992, time.April, 29, 0, 0, 0, 0, time.UTC),
		booksInUse:  []booksInUse{},
	},
}

type book struct {
	Title  string
	Copies int
	Author string
	Issue  time.Time
}

func postBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookByTitle(c *gin.Context) {
	title := c.Param("Title")

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
	booksInUse  []booksInUse
}

type booksInUse struct {
	NameOfBook book
	DateOfRent time.Time
}

func getReaderBooks(c *gin.Context) {
	var input struct {
		Name string `json:"Name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	readerName := input.Name

	for _, localReader := range readers {
		if localReader.Name == readerName {
			var okbooks []booksInUse
			var badbooks []booksInUse

			for _, rentedbook := range localReader.booksInUse {
				if time.Since(rentedbook.DateOfRent) <= timeToExpire {
					okbooks = append(okbooks, rentedbook)

				} else {
					badbooks = append(badbooks, rentedbook)
				}
			}

			c.IndentedJSON(http.StatusOK, gin.H{
				"okbooks":  okbooks,
				"badbooks": badbooks,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reader not found"})
}

func rentBookByTitle(c *gin.Context) {
	var input struct {
		Name  string `json:"Name"`
		Title string `json:"Title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	readerName := input.Name
	bookName := input.Title
	bookID := -1

	for id, localBook := range books {
		if localBook.Title == bookName {
			if localBook.Copies == 0 {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "All books are rented"})
				return
			}
			bookID = id
			break
		}
	}

	if bookID == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	for id, localReader := range readers {
		if localReader.Name == readerName {
			if len(localReader.booksInUse) >= 3 {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "You have too much books RN"})
				return
			}

			books[bookID].Copies--
			var newRent = booksInUse{books[bookID], time.Now()}

			readers[id].booksInUse = append(readers[id].booksInUse, newRent)
			c.IndentedJSON(http.StatusOK, gin.H{
				"reader": readers[id], "booksInUse": readers[id].booksInUse,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reader not found"})

}

func returnBookByTitle(c *gin.Context) {
	var input struct {
		Name  string `json:"Name"`
		Title string `json:"Title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	readerName := input.Name
	bookName := input.Title
	bookID := -1

	for id, localBook := range books {
		if localBook.Title == bookName {
			bookID = id
			break
		}
	}

	if bookID == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	for id, localReader := range readers {
		if localReader.Name == readerName {

			bookID := -1

			for id, book := range localReader.booksInUse {
				if book.NameOfBook.Title == bookName {
					bookID = id
					break
				}
			}
			if bookID == -1 {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book is not rented by you"})
				return
			}

			books[bookID].Copies++
			readers[id].booksInUse = append(readers[id].booksInUse[:bookID], readers[id].booksInUse[bookID+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{
				"reader": readers[id], "booksInUse": readers[id].booksInUse,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reader not found"})

}

func getReaders(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, readers)
}

func postReader(c *gin.Context) {
	var newReader book

	if err := c.BindJSON(&newReader); err != nil {
		return
	}

	books = append(books, newReader)
	c.IndentedJSON(http.StatusCreated, newReader)
}

func newYear(year int) time.Time {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func main() {
	router := gin.Default()
	router.GET("/readers", getReaders)
	router.POST("/reader", postReader)
	router.GET("/reader/books", getReaderBooks)

	router.GET("/books", getBooks)
	router.POST("/book", postBook)
	router.GET("/book/:Title", getBookByTitle)
	router.PATCH("/book/rent", rentBookByTitle)
	router.PATCH("/book/return", returnBookByTitle)

	router.Run("localhost:8080")
}
