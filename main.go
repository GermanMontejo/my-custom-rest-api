package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var myNotes = make(map[string]Note)
var id int

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon"`
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	idMap := mux.Vars(r)
	idStr := idMap["id"]
	value, ok := myNotes[idStr]

	var noteToUpdate Note
	err := json.NewDecoder(r.Body).Decode(&noteToUpdate)

	if err != nil {
		panic(err)
	}

	if !ok {
		fmt.Printf("Note entry with id: %s not found", idStr)
	} else {
		noteToUpdate.CreatedOn = value.CreatedOn
		delete(myNotes, idStr)
		myNotes[idStr] = noteToUpdate
	}
	w.WriteHeader(http.StatusOK)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var note Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}

	note.CreatedOn = time.Now()
	noteId := strconv.Itoa(id)
	myNotes[noteId] = note
	id++

	jsonNote, err := json.Marshal(note)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonNote)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	for _, value := range myNotes {
		notes = append(notes, value)
	}

	mNotes, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Notes: ", myNotes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(mNotes)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)
	idStr := ids["id"]
	_, ok := myNotes[idStr]
	if !ok {
		fmt.Printf("Failed to find note for:%s", idStr)
	} else {
		delete(myNotes, idStr)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	myMux := mux.NewRouter().StrictSlash(false)
	myMux.HandleFunc("/api/notes", PostHandler).Methods("POST")
	myMux.HandleFunc("/api/notes", GetHandler).Methods("GET")
	myMux.HandleFunc("/api/notes/{id}", PutHandler).Methods("PUT")
	myMux.HandleFunc("/api/notes/{id}", DeleteHandler).Methods("DELETE")
	server := &http.Server{
		Addr:    ":8080",
		Handler: myMux,
	}
	fmt.Println("Listening on port 8080.")
	server.ListenAndServe()
}
