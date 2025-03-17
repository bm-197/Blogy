package blog

import (
	"bytes"
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory sqlite database and initializes the posts table.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite database: %v", err)
	}
	// Create the posts table.
	sqlStmt := `
	CREATE TABLE posts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		t.Fatalf("Failed to create posts table: %v", err)
	}
	return db
}

// TestTemplateFile verifies that the template file exists and is not empty.
// This uses os.ReadFile, the modern replacement for ioutil.ReadFile.
func TestTemplateFile(t *testing.T) {
	// Adjust the path if tests are run from a different working directory.
	data, err := os.ReadFile("templates/post.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Template file is empty")
	}
}

// TestGetHandler checks that a GET request renders the blog post form.
func TestGetHandler(t *testing.T) {
	// Override the global template with a dummy template for testing.
	tmpl = template.Must(template.New("post").Parse(`<html><body><h1>Create a New Blog Post</h1></body></html>`))

	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/post", nil)
	w := httptest.NewRecorder()

	handler := PostBlogHandler(db)
	handler(w, req)

	res := w.Result()
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", res.StatusCode)
	}
	if !strings.Contains(string(body), "Create a New Blog Post") {
		t.Fatalf("Expected body to contain the form header, got: %s", string(body))
	}
}

// TestPostHandler verifies that a POST request creates a blog post.
func TestPostHandler(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Override the template to avoid dependency on an external file.
	tmpl = template.Must(template.New("post").Parse(`<html><body><h1>Create a New Blog Post</h1></body></html>`))

	// Create form data.
	data := url.Values{}
	data.Set("title", "Test Title")
	data.Set("content", "Test Content")

	req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler := PostBlogHandler(db)
	handler(w, req)

	res := w.Result()
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", res.StatusCode)
	}
	if !strings.Contains(string(body), "Blog post added successfully!") {
		t.Fatalf("Expected success message, got: %s", string(body))
	}

	// Verify that the post was added to the database.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query posts count: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 post in database, got %d", count)
	}
}

