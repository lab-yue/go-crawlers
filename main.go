package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// APIData is api raw response json
type APIData struct {
	Data Book `json:"data"`
}

// Book is "data" in response json
type Book struct {
	Volumes []Volume `json:"vs"`
}

// Volume is volume in a book
type Volume struct {
	Title    string    `json:"vN"`
	Chapters []Chapter `json:"cs"`
}

// Chapter is chapter in volume
type Chapter struct {
	Title string `json:"cN"`
	Time  string `json:"uT"`
	Words int    `json:"cnt"`
	URL   string `json:"cU"`
}

func main() {
	if len(os.Args) > 1 {
		book := GetBook(os.Args[1])
		fmt.Println(book)
	} else {
		fmt.Println("ERROR: NO BookID")
		os.Exit(2)
	}
}

// GetBook : get book data from api
func GetBook(BookID string) *Book {
	res := APIData{}
	req := &http.Client{Timeout: 10 * time.Second}
	url := "https://read.qidian.com/ajax/book/category?&bookId=" + BookID
	r, err := req.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&res)

	return &res.Data
}

func (book *Book) String() string {
	var lines []string

	for _, vol := range book.Volumes {
		lines = append(lines, vol.Title)
		for _, chapter := range vol.Chapters {
			lines = append(lines, "\t"+chapter.Title+" "+strconv.Itoa(chapter.Words)+"字")
		}
	}
	return strings.Join(lines, "\r\n")
}
