package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json: "id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUser struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	LoginAt  time.Time `json:"login_at"`
}

var err error

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

func NewMemberHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading .env file")
	}

	// 환경변수를 이용하여서 DB 접속 정보를 가지고 옴.
	DB_ACCOUNT := os.Getenv("DB_ACCOUNT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	Connection := DB_ACCOUNT + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME
	db, err := sql.Open("mysql", Connection)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	user := new(User)
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	UserID := uuid.Must(uuid.NewV4())
	if err != nil {
		panic(err)
	}

	pwHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.ID = UserID
	user.Password = string(pwHash)
	user.CreatedAt = time.Now()

	_, err = db.Exec("insert into Users (ID, FirstName, LastName, Email, Password, CreatedAt) value (?, ?, ?, ?, ?, ?)", user.ID, user.FirstName, user.LastName, user.Email, string(pwHash), user.CreatedAt)
	if err != nil {
		panic(err)
	}

	db.Close()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func LoginMemberHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading .env file")
	}

	// 환경변수를 이용하여서 DB 접속 정보를 가지고 옴.
	DB_ACCOUNT := os.Getenv("DB_ACCOUNT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	Connection := DB_ACCOUNT + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME
	db, err := sql.Open("mysql", Connection)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	var (
		Email    string
		Password string
	)

	LoginUser := new(LoginUser)
	err = json.NewDecoder(r.Body).Decode(&LoginUser)
	if err != nil {
		panic(err)
	}

	err = db.QueryRow("SELECT Email,Password FROM Users").Scan(&Email, &Password)
	if err != nil {
		panic(err)
	}

	LoginUser.Email = Email
	LoginUser.Password = Password
	LoginUser.LoginAt = time.Now()
}

func NewHandler() http.Handler {
	mux := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public/"))

	// GET | Render
	mux.HandleFunc("/", IndexRenderHandler).Methods("GET")
	mux.HandleFunc("/register", RegisterRenderHandler).Methods("GET")
	mux.HandleFunc("/login", LoginRenderHandler).Methods("GET")

	mux.HandleFunc("/register/new", NewMemberHandler).Methods("POST")
	mux.HandleFunc("/login", LoginMemberHandler).Methods("POST")

	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return mux
}
