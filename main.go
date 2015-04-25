package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinian/dice"
)

func rollHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.PostFormValue("user_name")
	desc := r.PostFormValue("text")
	result, _ := dice.Roll(desc)
	fmt.Fprintf(w, "%s rolled: %v", user, result)
}

func main() {
	http.HandleFunc("/roll", rollHandler)
	log.Fatal(http.ListenAndServe(":8989", nil))
}
