package main

import (
	"e-commerce_api/routing"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	routing.SetUp(router)

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, router))
}
