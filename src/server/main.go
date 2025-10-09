package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileserver := http.FileServer(http.Dir("../static"))

	http.Handle("/", fileserver)

	fmt.Printf("port running on http://localhost:8081/\n")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
