package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"user"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

func getUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	result := User{
		Username: params["id"],
		Commits:  0,
		Color:    0,
	}

	json.NewEncoder(writer).Encode(result)
}

func main() {
	fmt.Println("API Start:" + string(time.Now().Format("2006-01-02 15:04:05")))
	router := mux.NewRouter()

	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}
