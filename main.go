package main

import (
	"fiber-login/app"
	"log"
	"net/http"
)

func main() {
	log.Print("Listen And Server at http://localhost:3000")
	http.ListenAndServe(":3000", app.NewHandler())
}
