package main

import (
	"fmt"
	"net/http"
	"os"
	//"strings"
)

func main() {
	fmt.Println("PORT: " + os.Getenv("PORT"))

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Http Server Test!")
		fmt.Println(r.Header)
	})

	/*
		http.ListenAndServe(os.Getenv("PORT"), nil)
	*/
}
