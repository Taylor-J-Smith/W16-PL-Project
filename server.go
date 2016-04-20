package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

// Basic Page struct.
type Page struct {
	Title string
	Body  []byte
}

// Method that is called on a page to save it to file.
// Filename is the title of the page.
// If file does not exist it is created with permisions 0600 (rw for user)
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Loads a page from memory with the title given in the parameter.
// Returns a nil string if the page does not exist otherwise returns a page.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

  if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

// Basic handler used for demo purposes.
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
