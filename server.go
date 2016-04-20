package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    "regexp"
    "errors"
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
	filename := "pages/" + title + ".txt"

  body, err := ioutil.ReadFile(filename)

  if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {

    m := validPath.FindStringSubmatch(r.URL.Path)

    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }

    return m[2], nil // The title is the second subexpression.
}

// Basic handler used for demo purposes.
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Redirects the root call to the front page of the wiki.
func rootRedirectHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

    title, err := getTitle(w, r)

    if err != nil {
        return
    }

    p, err := loadPage(title)

    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)

    p, err := loadPage(title)

    if err != nil {
        p = &Page{Title: title}
    }

    renderTemplate(w, "edit", p)
}

func main() {

    http.HandleFunc("/", rootRedirectHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.ListenAndServe(":8080", nil)
}
