package app

import (
	"database/sql"
	"fmt"
	"go-login/render"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
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
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	LoginAt   time.Time
}

var (
	err     error
	store   = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	errCode http.ConnState
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
	err = godotenv.Load(".env")
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

	UserPostEmail := r.PostFormValue("email")
	UserPostPassword := r.PostFormValue("password")

	rows, err := db.Query("SELECT  Email, ID, Password FROM Users where Email= ?", UserPostEmail)
	if err != nil {
		fmt.Printf("DB ERROR")
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&User.Email, &User.ID, &User.Password)
		User.LoginAt = time.Now()
		fmt.Println(User)
		err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(UserPostPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Print("BAD ERROR")
		} else {
			w.WriteHeader(http.StatusOK)
			session, _ := store.Get(r, "auth-login")
			// IP를 기록하는 코드를 작성하였지만 되지 않아서 추후 개발하여서 추가할 예정임. / User.IP Noting
			session.Values["ID"] = User.ID
			session.Values["Email"] = User.Email
			session.Values["LogindAt"] = User.LoginAt
			session.Save(r, w)
			fmt.Println(User)

			// ERROR: 1회만 되고 그 후 업데이트가 안 됨.
			_, err = db.Exec("insert into autu_login (ID, Email, LogindAt) value (?, ?, ?)", User.ID, User.Email, User.LoginAt)
			if err != nil {
				panic("ERROR EXEC SESSION LOG\n")
			}
			h := `<!DOCTYPE html>
			<html>
			<script>
				alert("로그인 성공!");
				location.herf="/public/view/home/index.html";
			</script>
			</html>`
			fmt.Fprint(w, h)
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
