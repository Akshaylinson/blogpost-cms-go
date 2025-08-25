package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

// -----------------
// Data Structures
// -----------------
type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
}

var (
	posts     []Post
	postID    int
	mutex     sync.Mutex
	templates = template.Must(template.ParseGlob("templates/*.html"))
)

// -----------------
// Handlers
// -----------------
func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", posts)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")

		mutex.Lock()
		postID++
		newPost := Post{ID: postID, Title: title, Content: content, Author: author}
		posts = append(posts, newPost)
		mutex.Unlock()

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "admin.html", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var selectedPost Post
	for _, p := range posts {
		if p.ID == id {
			selectedPost = p
			break
		}
	}

	templates.ExecuteTemplate(w, "post.html", selectedPost)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/post", postHandler)

	println("🚀 Blog CMS running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
