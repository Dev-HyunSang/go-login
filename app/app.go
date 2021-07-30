package app

import (
	"database/sql"
	"fmt"
	"go-login/render"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"
)

/*Sturct를 두 개로 분할하였습니다. 이유는 로그인과 회원가입시 Struct 더 확실히 구분하기 위함입니다.
User = NewMemberHandler에서 사용되며 회원가입시 DB로 구조화 시켜서 Insert함.
LoginUser = LoginMemberHandler에서 사용하고 있으며 로그인시 DB로 SELETC를 하기 위함. */
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type LoginUser struct {
	IP        net.IP
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	LoginAt   time.Time
}

var (
	err   error
	key   = []byte("super-secret-key")
	store *mysqlstore.MySQLStore
)

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

	UserID := uuid.Must(uuid.NewV4())
	if err != nil {
		panic(err)
	}
	user.Password = r.FormValue("password")

	pwHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.ID = UserID
	user.FirstName = r.FormValue("first_name")
	user.LastName = r.FormValue("last_name")
	user.Email = r.FormValue("email")
	user.Password = string(pwHash)
	user.CreatedAt = time.Now()

	_, err = db.Exec("insert into Users (ID, FirstName, LastName, Email, Password, CreatedAt) value (?, ?, ?, ?, ?, ?)", user.ID, user.FirstName, user.LastName, user.Email, string(pwHash), user.CreatedAt)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
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

	User := new(LoginUser)

	var (
		Email    string
		ID       uuid.UUID
		Password string
	)

	UserPostEmail := r.PostFormValue("email")
	UserPostPassword := r.PostFormValue("password")

	rows, err := db.Query("SELECT  Email, ID, Password FROM Users where Email= ?", UserPostEmail)
	if err != nil {
		fmt.Printf("DB ERROR")
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&Email, &ID, &Password)
		User.Email = Email
		User.ID = ID
		User.Password = Password
		err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(UserPostPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			// GET LOGIN UERES IP
			addres, err := net.InterfaceAddrs()
			if err != nil {
				os.Stderr.WriteString("Oops: " + err.Error() + "\n")
				os.Exit(1)
			}

			for _, a := range addres {
				if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						User.IP = ipnet.IP
					}
				}
			}

			w.WriteHeader(http.StatusOK)
			sessions, err := store.Get(r, "login-sessino")
			sessions.Values["IP"] = User.IP
			sessions.Values["UUID"] = User.ID
			sessions.Values["Email"] = User.Email
			User.LoginAt = time.Now()
			sessions.Values["LoginAt"] = User.LoginAt
			err = sessions.Save(r, w)
		}
	}
}

func NewHandler() http.Handler {
	mux := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public/"))

	// GET | Render
	mux.HandleFunc("/", render.IndexRenderHandler).Methods("GET")
	mux.HandleFunc("/register", render.RegisterRenderHandler).Methods("GET")
	mux.HandleFunc("/login", render.LoginRenderHandler).Methods("GET")
	mux.HandleFunc("/home/index", render.HomeRenderHandler).Methods("GET")

	mux.HandleFunc("/register", NewMemberHandler).Methods("POST")
	mux.HandleFunc("/login", LoginMemberHandler).Methods("POST")

	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return mux
}
