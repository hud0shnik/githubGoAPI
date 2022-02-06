package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		return User{}
	}

	// Запись информации из респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		// Получение сегодняшней даты
		// Добавляет 3 часа т.к сервер находится в другом часовом поясе
		date = string(time.Now().Add(time.Hour * 3).Format("2006-01-02"))
	}

	// Вот так выглядит html одной ячейки:
	// <rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1">

	// Поиск сегодняшней ячейеки
	if strings.Contains(string(body), "data-date=\""+date) {
		pageStr, i := string(body), 14000

		for ; i < len(pageStr)-24000; i++ {
			if pageStr[i:i+21] == "data-date=\""+date {
				for ; pageStr[i] != '<'; i-- {
					// Доводит i до начала кода ячейки
				}
				break
			}
		}
		// Получение параметров ячейки
		values := strings.FieldsFunc(pageStr[i:i+155], func(r rune) bool {
			return r == '"'
		})

		// Запись и обработка нужной информации
		dataLevel, _ := strconv.Atoi(values[15])
		commits, _ := strconv.Atoi(values[19])

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

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getCommits(mux.Vars(request)["id"], mux.Vars(request)["date"]))
}

func main() {
	// Вывод времени начала работы
	fmt.Println("API Start:" + string(time.Now().Add(time.Hour*3).Format("2006-01-02 15:04:05")))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/{id}", getInfo).Methods("GET")
	router.HandleFunc("/{id}/{date}", getInfo).Methods("GET")

	// Запуск API
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
