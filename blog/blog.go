package blog

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Post represents a blog post.
type Post struct {
	ID      int
	Title   string
	Content string
}

// Parse the template used for the blog posting form.
// Ensure the "templates" directory exists at the project root with "post.html".
var tmpl = template.Must(template.ParseFiles("templates/post.html"))

// PostBlogHandler handles GET and POST requests for blog posts.
func PostBlogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

