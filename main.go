package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bm-197/Blogy/blog"
)

func main() {
	// Open the SQLite database (or create if it doesn't exist)
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the posts table if it does not exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT
	);
	`
	if _, err = db.Exec(sqlStmt); err != nil {
		log.Fatalf("Failed to create table: %q", err)
	}

	// Register the blog post handler
	http.HandleFunc("/post", blog.PostBlogHandler(db))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

