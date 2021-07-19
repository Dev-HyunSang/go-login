package app

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID        int       `json: "id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	lastID int
)

// Only Render Handler and Method "GET"
func IndexRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func RegisterRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/register.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func LoginRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/login.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func NewHandler() http.Handler {
	mux := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public/"))

	// GET | Render
	mux.HandleFunc("/", IndexRenderHandler).Methods("GET")
	mux.HandleFunc("/register", RegisterRenderHandler).Methods("GET")
	mux.HandleFunc("/login", LoginRenderHandler).Methods("GET")

	// POST

	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return mux
}
