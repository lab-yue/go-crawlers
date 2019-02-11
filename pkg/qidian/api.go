package qidian

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
	c := make(chan *Book)
	go GetBook(c, os.Args[1])
	fmt.Println("waiting network")
	if len(os.Args) > 1 {
		for {
			select {
			case book := <-c:
				fmt.Println("\n\r", book)
				return
			default:
				fmt.Print(".")
				time.Sleep(25 * time.Millisecond)
			}
		}
	} else {
		fmt.Println("ERROR: NO BookID")
		os.Exit(2)
	}
}

// GetBook : get book data from api
func GetBook(c chan *Book, BookID string) {
	res := APIData{}
	req := &http.Client{Timeout: 10 * time.Second}
	url := "https://read.qidian.com/ajax/book/category?&bookId=" + BookID
	r, err := req.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&res)

	c <- &res.Data
	return
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
