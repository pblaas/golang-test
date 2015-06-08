package main

import (
	"net/http"
)

func main() {
	http.Handle("/submit", http.FileServer(http.Dir("public/submit")))
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":3000", nil)
}
