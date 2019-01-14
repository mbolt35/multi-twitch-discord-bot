package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Go Http Server Test!")
		fmt.Println(r.Header)
	})

	http.ListenAndServe(os.Genenv("PORT"), nil)
}
