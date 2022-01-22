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
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

func getCommits(username string, date string) User {
	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		fmt.Println("Github error: ", err)
		return User{}
	}

	// Запись информации из респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		// Получение сегодняшней даты
		date = string(time.Now().Add(time.Hour * 3).Format("2006-01-02"))
	}

	// Вот так выглядит html одной ячейки:
	// <rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1"></rect>

	// Поиск сегодняшней ячейеки
	if strings.Contains(string(body), "data-date=\""+date) {
		pageStr, commitsString, i := string(body), "", 0

		// Проход по всему html файлу в поисках нужной клетки
		for ; i < len(pageStr)-40; i++ {
			if pageStr[i:i+21] == "data-date=\""+date {
				// Так как количество коммитов стоит перед датой, переставляем i
				i -= 7
				break
			}
		}
		for ; pageStr[i] != '"'; i++ {
			// Доводит i до символа "
		}
		for i++; pageStr[i] != '"'; i++ {
			// Считывание и запись значения в скобках
			commitsString += string(pageStr[i])
		}
		for i += 35; pageStr[i] != '"'; i++ {
		}

		// Запись и обработка полученной информации
		dataLevel, _ := strconv.Atoi(pageStr[i+1 : i+2])
		commits, _ := strconv.Atoi(commitsString)

		return User{
			Date:     date,
			Username: username,
			Commits:  commits,
			Color:    dataLevel,
		}
	}

	return User{
		Date:     date,
		Username: username,
		Commits:  0,
		Color:    0,
	}
}

func getInfo(writer http.ResponseWriter, request *http.Request) {
	// Заголовок, определяющий тип данных для работы
	writer.Header().Set("Content-Type", "application/json")

	// Запись параметров из реквеста
	params := mux.Vars(request)
	// Обработка данных и вывод результата

	result := getCommits(params["id"], params["date"])
	json.NewEncoder(writer).Encode(result)
}

func main() {
	// Вывод даты начала работы
	port := os.Getenv("PORT")
	fmt.Println("API Start:" + string(time.Now().Add(time.Hour*3).Format("2006-01-02 15:04:05")))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/user/{id}", getInfo).Methods("GET")
	router.HandleFunc("/user/{id}/{date}", getInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
