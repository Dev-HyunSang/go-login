package main

import (
	"go-login/app"
	"log"
	"net/http"
)

func main() {
	log.Print("Listen And Server at http://localhost:3000")
	err := http.ListenAndServe(":3000", app.NewHandler())
	if err != nil {
		panic(err)
	}
}
