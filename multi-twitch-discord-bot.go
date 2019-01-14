package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0] + os.Getenv(pair[0]))
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Println("$PORT not set. Defaulting to 3001")
		port = "3001"
	}

	users := [3]string{"bolt_", "lucky918", "perlj"}

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe(":"+port, nil)
}
