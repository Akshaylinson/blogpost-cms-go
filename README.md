BlogCMS (Go + HTML Mini CMS)

A lightweight Content Management System (CMS) built using Go (Golang) and plain HTML templates.
This project demonstrates how an Admin can add posts (title, content, images, and links) via an admin panel, and how Users can view those posts on the website.

🚀 Features

Admin Panel

Add new posts with title, content, images, and URLs

Manage multiple posts dynamically

User Side

View published posts in a clean layout

Posts automatically update when added by admin

Tech Stack

Go (Golang) – for backend and server

HTML/CSS – for frontend templates

net/http + html/template – Go standard libraries (no extra dependencies)

📂 Project Structure
blogcms/
│── main.go          # Go server and routes
│── index.html       # User landing page
│── admin.html       # Admin dashboard to add posts
│── post.html        # Template for individual posts
│── go.mod           # Go module file

⚙️ Setup & Run

Clone the repo:

git clone https://github.com/your-username/blogcms.git
cd blogcms


Initialize Go module:

go mod init blogcms
go mod tidy


Run the server:

go run main.go


Open in browser:

User Page → http://localhost:8080

Admin Page → http://localhost:8080/admin

🛠️ Future Improvements

Add authentication for admin panel

Store posts in a database (SQLite / PostgreSQL) instead of memory

Add images upload support

Improve frontend design with Bootstrap/Tailwind

📜 License

This project is open-source and available under the MIT License.
