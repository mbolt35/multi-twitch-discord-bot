package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Go Http Server Test!")
		fmt.Println(r.Header)
	})

	http.ListenAndServe(":3001", nil)
}
