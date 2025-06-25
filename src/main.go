package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	//mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /scribble", postScribble)
	mux.HandleFunc("GET /scribble", getScribble)

	fmt.Println("Server Listening to :8080")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "well hello there, it's the root.")
}

func getScribble(
	w http.ResponseWriter,
	r *http.Request,
) {
	scribble, enough := randScribble()
	if !enough {
		http.Error(w, "there are no current scribbles", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(*scribble)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func postScribble(
	w http.ResponseWriter,
	r *http.Request,
) {
	var scribble Scribble
	err := json.NewDecoder(r.Body).Decode(&scribble)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if scribble.Text == "" {
		http.Error(w, "scribble text can not be empty", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(scribble.Text) > 256 {
		http.Error(w, "scribble text can not be longer than 256 characters", http.StatusBadRequest)
		return
	}

	addScribble(scribble)

	w.WriteHeader(http.StatusNoContent)
}
