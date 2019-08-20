package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// STRUCTS or MODELS
type Book struct {
	ID     string  `json:"id"` // Syntax to parse json to struct
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"` // "*" in front of variable = value of variable (pointer)
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// INIT BOOK VARIABLE AS A SLICE BOOK STRUCT
var books []Book // a Slice is data structure similar to arrays, but dynamic

// ENDPOINTS

func getBooks(w http.ResponseWriter, r *http.Request) { // "*" in front of type = type to what tha variable points
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books) //Encode() receives the element we want to output
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{}) // Return empty book if id doesn't exist
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)  // "_" stores error returned by Decode()
	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID, don't use in production!
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Get rid of the book using JS slice() style
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Get rid of the book using JS slice() style
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {

	// Init mux router
	router := mux.NewRouter() //":=" syntax allows type inference

	// Mock data. @todo: implement DB
	books = append(books, Book{
		ID:    "1",
		Isbn:  "123-ABC",
		Title: "Book One",
		Author: &Author{ // "&" in front of variable: variable's memory address
			Firstname: "John",
			Lastname:  "Doe",
		},
	})

	books = append(books, Book{
		ID:    "2",
		Isbn:  "456-DEF",
		Title: "Book Two",
		Author: &Author{
			Firstname: "Steve",
			Lastname:  "Smith",
		},
	})

	books = append(books, Book{
		ID:    "3",
		Isbn:  "789-GHI",
		Title: "Book Three",
		Author: &Author{
			Firstname: "Michael",
			Lastname:  "Johnson",
		},
	})

	// Route/Endpoint handlers
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Run server
	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8000", router)) // Store error in fatal log
}
