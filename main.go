package main

import (
	"coupon/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	log.Println("Server is starting... Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
