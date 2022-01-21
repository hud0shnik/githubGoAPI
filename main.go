package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"user"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

func getCommits(username string) User {
	fmt.Println("Checking user: " + username)
	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		fmt.Println("Github error: ", err)
		return User{}
	}
	// Запись информации из респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	currentDate := string(time.Now().Format("2006-01-02"))

	// Вот так выглядит html одной ячейки:
	// <rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1"></rect>
	if strings.Contains(string(body), "data-date=\""+currentDate+"\" data-level=\"") {
		pageStr, commitsString := string(body), ""
		i := 0

		// Проход по всему html файлу в поисках нужной клетки
		for ; i < len(pageStr)-40; i++ {
			if pageStr[i:i+35] == "data-date=\""+currentDate+"\" data-level=\"" {
				// Так как количество коммитов стоит перед датой, переставляем i
				i -= 7
				break
			}
		}
		for ; pageStr[i] != '"'; i++ {
			// Доводит i до символа "
		}
		for i++; pageStr[i] != '"'; i++ {
			// Считывание значения в скобках
			commitsString += string(pageStr[i])
		}
		for i += 35; pageStr[i] != '"'; i++ {
		}
		// Запись и обработка полученной информации (цвет клетки)
		dataLevel, _ := strconv.Atoi(pageStr[i+1 : i+2])
		commits, _ := strconv.Atoi(commitsString)
		return User{
			Username: username,
			Commits:  commits,
			Color:    dataLevel,
		}
	}
	return User{
		Username: username,
		Commits:  0,
		Color:    0,
	}
}

func getUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	result := getCommits(params["id"])

	json.NewEncoder(writer).Encode(result)
}

func main() {
	fmt.Println("API Start:" + string(time.Now().Format("2006-01-02 15:04:05")))
	router := mux.NewRouter()

	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}
