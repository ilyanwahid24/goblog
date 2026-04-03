package handlers

import (
	"context"
	"go-blog/models"
	"html/template"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	UserStore    *models.UserStore
	SessionStore *models.SessionStore
	Templates    map[string]*template.Template
}

func NewAuthHandler(userStore *models.UserStore, sessionStore *models.SessionStore) *AuthHandler {
	funcMap := template.FuncMap{}

	templates := make(map[string]*template.Template)
	pages := []string{"login.html"}
	for _, page := range pages {
		t := template.Must(
			template.New("").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+page),
		)
		templates[page] = t
	}

	return &AuthHandler{
		UserStore:    userStore,
		SessionStore: sessionStore,
		Templates:    templates,
	}
}

// LoginPage renders the login form
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}

	if err := h.Templates["login.html"].ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error rendering login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// LoginPost processes the login form submission
func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.UserStore.Authenticate(username, password)
	if err != nil {
		data := map[string]interface{}{
			"Title": "Login",
			"Error": "Invalid username or password",
		}
		h.Templates["login.html"].ExecuteTemplate(w, "layout", data)
		return
	}

	// Create session
	token, err := h.SessionStore.CreateSession(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// Logout deletes the session
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.SessionStore.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AuthMiddleware protects routes
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := h.SessionStore.ValidateSession(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-1 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
