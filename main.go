// main.go
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// -----------------
// Enhanced Data Structures
// -----------------
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Excerpt   string    `json:"excerpt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Published bool      `json:"published"`
	Tags      []string  `json:"tags"`
	ImageURL  string    `json:"image_url"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // In production, use hashed passwords
	Role     string `json:"role"`     // "admin" or "editor"
}

type BlogData struct {
	Posts       []Post
	Categories  []string
	RecentPosts []Post
	CurrentPost *Post
}

var (
	posts      []Post
	users      []User
	postID     int
	mutex      sync.Mutex
	templates  *template.Template
	categories = []string{"Technology", "Lifestyle", "Travel", "Food", "Business"}
)

// -----------------
// Initialization
// -----------------
func init() {
	// Sample data for demonstration
	posts = []Post{
		{
			ID:        1,
			Title:     "Welcome to Our Professional Blog",
			Content:   "This is a sample blog post to get things started. Our new CMS provides a clean, modern interface for both readers and content creators. With features like categories, tags, and a responsive design, we're ready to deliver great content to our audience.",
			Excerpt:   "A warm welcome to our new blog platform with enhanced features",
			Author:    "Admin",
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now().Add(-72 * time.Hour),
			Published: true,
			Tags:      []string{"welcome", "blog", "update"},
			ImageURL:  "/static/images/sample-blog.jpg",
		},
		{
			ID:        2,
			Title:     "The Future of Web Development",
			Content:   "Web development continues to evolve at a rapid pace. With new frameworks, tools, and methodologies emerging regularly, developers must stay current to remain competitive. In this post, we explore the latest trends and what they mean for the future of web development.",
			Excerpt:   "Exploring the latest trends in web development and what's coming next",
			Author:    "Jane Developer",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-48 * time.Hour),
			Published: true,
			Tags:      []string{"webdev", "technology", "programming"},
			ImageURL:  "/static/images/webdev.jpg",
		},
	}
	postID = 2

	// Sample users
	users = []User{
		{Username: "admin", Password: "admin123", Role: "admin"},
		{Username: "editor", Password: "editor123", Role: "editor"},
	}

	// Parse templates with functions
	templateFuncs := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"truncate": func(s string, length int) string {
			if len(s) < length {
				return s
			}
			return s[:length] + "..."
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	templates = template.Must(template.New("").Funcs(templateFuncs).ParseGlob("templates/*.html"))
}

// -----------------
// Middleware
// -----------------
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("APP_ENV") == "development" {
			// Skip authentication in development for easier testing
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check credentials
		authenticated := false
		for _, user := range users {
			if user.Username == username && user.Password == password {
				authenticated = true
				break
			}
		}

		if !authenticated {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// -----------------
// Handlers
// -----------------
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Get only published posts for public view
	var publishedPosts []Post
	for _, post := range posts {
		if post.Published {
			publishedPosts = append(publishedPosts, post)
		}
	}

	// Get recent posts (last 5)
	var recentPosts []Post
	if len(publishedPosts) > 5 {
		recentPosts = publishedPosts[:5]
	} else {
		recentPosts = publishedPosts
	}

	data := BlogData{
		Posts:       publishedPosts,
		Categories:  categories,
		RecentPosts: recentPosts,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		excerpt := r.FormValue("excerpt")
		tags := r.Form["tags"]
		published := r.FormValue("published") == "on"
		imageURL := r.FormValue("image_url")

		mutex.Lock()
		postID++
		newPost := Post{
			ID:        postID,
			Title:     title,
			Content:   content,
			Excerpt:   excerpt,
			Author:    author,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Published: published,
			Tags:      tags,
			ImageURL:  imageURL,
		}
		posts = append(posts, newPost)
		mutex.Unlock()

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "admin.html", map[string]interface{}{
		"Posts":      posts,
		"Categories": categories,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var selectedPost Post
	for _, p := range posts {
		if p.ID == id {
			selectedPost = p
			break
		}
	}

	if selectedPost.ID == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Get recent posts for sidebar
	var recentPosts []Post
	for _, p := range posts {
		if p.Published && p.ID != selectedPost.ID {
			recentPosts = append(recentPosts, p)
			if len(recentPosts) >= 5 {
				break
			}
		}
	}

	data := BlogData{
		Posts:       posts,
		RecentPosts: recentPosts,
		Categories:  categories,
		CurrentPost: &selectedPost,
	}

	err = templates.ExecuteTemplate(w, "post.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Return only published posts for public API
	var publishedPosts []Post
	for _, post := range posts {
		if post.Published {
			publishedPosts = append(publishedPosts, post)
		}
	}

	err := json.NewEncoder(w).Encode(publishedPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "ok", "message": "Blog CMS is running"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up environment
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/health", healthCheckHandler)

	// Admin routes with authentication
	http.HandleFunc("/admin", authMiddleware(adminHandler))

	// API routes
	http.HandleFunc("/api/posts", apiPostsHandler)

	// Start server
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	log.Printf("🚀 Professional Blog CMS running in %s mode at http://localhost%s", os.Getenv("APP_ENV"), port)
	log.Printf("📊 Admin panel available at http://localhost%s/admin", port)
	log.Printf("❤️  Health check at http://localhost%s/health", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
