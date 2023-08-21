package main

import (
	"fmt"
	"gii/gii"
	"net/http"
)

func main() {
	Router := gii.New()

	Router.Get("/ping", handle)

	Router.Run("localhost:8000")
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "url path: %s\n", r.URL.Path)
}
