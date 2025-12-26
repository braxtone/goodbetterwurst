package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

//go:embed templates/**/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates = template.Must(
	template.ParseFS(templatesFS, "templates/**/*.html"),
)

func renderTemplate(w http.ResponseWriter, name string, p *Page) {
	err := templates.ExecuteTemplate(w, name, p)

	if err != nil {
		log.Printf("template execution failed (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

// func staticHandler(w http.ResponseWriter, r *http.Request) {
// 	trimmedPath := strings.Trim(r.URL.Path, "/")
// 	// Split the path into parts
// 	parts := strings.Split(trimmedPath, "/")
// 	n := len(parts)
// 	if n >= 2 {
// 		filePath := filepath.Join("./static", parts[n-2], parts[n-1])
// 		http.ServeFile(w, r, filePath)
// 	} else {
// 		http.Error(w, "Static file not found", http.StatusNotFound)
// 	}
// }
//

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// TODO:
// * Implement a module system so that we know, for the given user, which modules they can access
// * Based on that, generate the list of links that are live
// * Template the template so everything inherits from the same base template
func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "GoodBetterWurst", Body: []byte("")}
	renderTemplate(w, "index", p)
}

func percentageChartHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title: "Percentage Charts!",
		Body:  []byte("Coming Soon!")}

	renderTemplate(w, "pctcharts", p)
}

func main() {
	var fsys fs.FS
	var err error

	fsys, err = fs.Sub(staticFS, "static") // strip the top-level static/ folder
	if err != nil {
		log.Fatalf("failed to create sub filesystem: %v", err)
	}

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/percentageCharts/", percentageChartHandler)
	// http.HandleFunc("/static/", staticHandler)
	// TODO: look at adding caching for static images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
