package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0])
	}
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Http Server Test!")
		fmt.Println(r.Header)
	})
	/*


		http.ListenAndServe(os.Getenv("PORT"), nil)
	*/
}
