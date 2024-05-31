package main

import (
  "html/template"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strings"
)

type Page struct {
  Title string
  Body []byte
}

var tmpl *template.Template

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t, err := template.ParseFiles("./templates/" + tmpl + ".html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  t.Execute(w, p)
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return os.WriteFile(filename, p.Body, 0600)
}

// load a page by name "title" and return a pointer to a Page obj
func loadPage(title string) (*Page, error) {
  // TODO: fix the vulnerability here, path traversal
  filename := title + ".txt"
  // return the file contents to body and errs to _
  // is an `identifier` that is used as a throwaway
  body, err := os.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  // New object of type Page?
  return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, _ := loadPage(title)
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
  trimmedPath := strings.Trim(r.URL.Path, "/")
  // Split the path into parts
  parts := strings.Split(trimmedPath, "/")
  n := len(parts)
  if n >= 2 {
    filePath := filepath.Join("./static", parts[n-2], parts[n-1])
    http.ServeFile(w, r, filePath)
  } else {
    http.Error(w, "Static file not found", http.StatusNotFound)
  }
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  p.save()
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// TODO:
// * Implement a module system so that we know, for the given user, which modules they can access
// * Based on that, generate the list of links that are live
// * Template the template so everything inherits from the same base template
func indexHandler(w http.ResponseWriter, r *http.Request) {
  p := &Page{Title: "GoodBetterWurst", Body: []byte("Coming Soon!")}
  renderTemplate(w, "index", p)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.HandleFunc("/static/", staticHandler)
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
