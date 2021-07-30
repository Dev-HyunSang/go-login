package render

import (
	"net/http"
	"text/template"
)

func IndexRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/view/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func RegisterRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/view/register.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func LoginRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/view/login.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func HomeRenderHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/view/home/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}
