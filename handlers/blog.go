package handlers

import (
	"fmt"
	"go-blog/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	Store     *models.PostStore
	Templates map[string]*template.Template
}

func NewBlogHandler(store *models.PostStore) *BlogHandler {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"formatDate": func(t interface{}) string {
			switch v := t.(type) {
			case interface{ Format(string) string }:
				return v.Format("Jan 2, 2006")
			default:
				return fmt.Sprintf("%v", v)
			}
		},
		"truncate": func(s string, n int) string {
			if len(s) <= n {
				return s
			}
			return s[:n] + "..."
		},
	}

	templates := make(map[string]*template.Template)
	pages := []string{"home.html", "post.html", "admin.html", "edit.html"}
	for _, page := range pages {
		t := template.Must(
			template.New("").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+page),
		)
		templates[page] = t
	}

	return &BlogHandler{
		Store:     store,
		Templates: templates,
	}
}

// HomePage renders the blog homepage with all published posts
func (h *BlogHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Store.GetPublishedPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Home",
		"Posts": posts,
	}

	if err := h.Templates["home.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// PostPage renders a single blog post
func (h *BlogHandler) PostPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	post, err := h.Store.GetPostBySlug(slug)
	if err != nil {
		log.Printf("Post not found: %s - %v", slug, err)
		http.Error(w, "Post Not Found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": post.Title,
		"Post":  post,
	}

	if err := h.Templates["post.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AdminPage renders the admin panel
func (h *BlogHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Store.GetAllPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Admin",
		"Posts": posts,
	}

	if err := h.Templates["admin.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AdminEditPage renders the edit form for a post
func (h *BlogHandler) AdminEditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.Store.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post Not Found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Edit Post",
		"Post":  post,
	}

	if err := h.Templates["edit.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AdminNewPage renders the new post form
func (h *BlogHandler) AdminNewPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "New Post",
		"Post":  &models.Post{},
	}

	if err := h.Templates["edit.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreatePost handles the creation of a new post
func (h *BlogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post := &models.Post{
		Title:     strings.TrimSpace(r.FormValue("title")),
		Slug:      strings.TrimSpace(r.FormValue("slug")),
		Excerpt:   strings.TrimSpace(r.FormValue("excerpt")),
		Content:   r.FormValue("content"),
		Published: r.FormValue("published") == "on",
	}

	if post.Title == "" || post.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	if err := h.Store.CreatePost(post); err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// UpdatePost handles updating an existing post
func (h *BlogHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post := &models.Post{
		ID:        id,
		Title:     strings.TrimSpace(r.FormValue("title")),
		Slug:      strings.TrimSpace(r.FormValue("slug")),
		Excerpt:   strings.TrimSpace(r.FormValue("excerpt")),
		Content:   r.FormValue("content"),
		Published: r.FormValue("published") == "on",
	}

	if post.Title == "" || post.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	if err := h.Store.UpdatePost(post); err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// DeletePost handles deleting a post
func (h *BlogHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeletePost(id); err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
