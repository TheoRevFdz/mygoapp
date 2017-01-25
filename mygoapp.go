package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/css/", fs)
	http.HandleFunc("/show/", showHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", welcomeHandler)
	log.Println("Escuchando en http://localhost:8080/show/")
	http.ListenAndServe(":8080", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	page.save()
	http.Redirect(w, r, "/show/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	// template file
	t, err := template.ParseFiles("./views/edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, page)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bienvenidos</h1>")
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/show/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Fprintf(w, "<h1>%s</h1>", err)
	}

	// template file
	t, err := template.ParseFiles("./views/show.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, p)
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
