package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Println("Escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Fprintf(w, "<h1>%s</h1>", err)
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

// Page sirve para estructurar la pagina
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	err := ioutil.WriteFile(filename, p.Body, 0600)
	return err
}

func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	page := &Page{Title: title, Body: body}
	return page, err
}
