package main

import (
	"net/http"

	"github.com/bookstore/server"
)

func main() {
	r := server.StartService()
	http.ListenAndServe(":3000", r)
}
