package blog

import (
	"database/sql"
    "path/filepath"
	"html/template"
	"log"
	"net/http"
    "os"
)

// Post represents a blog post.
type Post struct {
	ID      int
	Title   string
	Content string
}


// getTemplate returns the parsed template.
func getTemplate() *template.Template {
	path := os.Getenv("TEMPLATES_PATH")
	if path == "" {
		path = "templates/post.html"
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to get absolute path for template: %v", err)
	}
	return template.Must(template.ParseFiles(absPath))
}

// PostBlogHandler handles GET and POST requests for blog posts.
func PostBlogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        tmpl := getTemplate()
		switch r.Method {
		case "GET":
			// Render the form for creating a new blog post.
			if err := tmpl.Execute(w, nil); err != nil {
				log.Printf("Template execution error: %v", err)
				http.Error(w, "Template error", http.StatusInternalServerError)
			}
		case "POST":
			// Parse the form data.
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
			title := r.FormValue("title")
			content := r.FormValue("content")
			if title == "" || content == "" {
				http.Error(w, "Title and Content are required", http.StatusBadRequest)
				return
			}
			// Insert the new post into the database.
			_, err := db.Exec("INSERT INTO posts(title, content) VALUES(?, ?)", title, content)
			if err != nil {
				log.Printf("DB insert error: %v", err)
				http.Error(w, "Failed to post blog", http.StatusInternalServerError)
				return
			}
			// Return a confirmation message (this works well with HTMX).
			w.Write([]byte("Blog post added successfully!"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

