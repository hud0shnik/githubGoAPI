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

// Структура для храния информации о пользователе
type User struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения информации с сайта
func getCommits(username string, date string) User {
	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return User{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Результат
	result := User{
		Date:     date,
		Username: username,
	}

	// HTML страницы в формате string
	pageStr := string(body[100000:225000])

	// Индекс символа начала ссылки
	left := strings.Index(pageStr, "https://avatars.githubusercontent.com/u")

	// Если ссылка найдена, считывает её и записывает
	if left != -1 {
		right := left
		for ; pageStr[right] != '?'; right++ {
			// Доводит pageStr[right] до символа '?'
		}

		// Запись ссылки
		result.Avatar = pageStr[left:right]
	}

	// Так выглядит html одной ячейки календаря:
	// <rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1">

	// Обрезает ненужную часть страницы
	pageStr = pageStr[50000:]

	// Указатель на ячейку нужной даты
	i := strings.Index(pageStr, "data-date=\""+date)

	// Проверка на существование нужной ячейки
	if i != -1 {
		for ; pageStr[i] != 's'; i-- {
			// Доводит i до     v
			//				class="ContributionCalendar-day"
		}

		// Получение параметров ячейки
		values := strings.FieldsFunc(pageStr[i:i+150], func(r rune) bool {
			return r == '"'
		})

		// Запись и обработка нужной информации
		result.Color, _ = strconv.Atoi(values[11])
		result.Commits, _ = strconv.Atoi(values[7])

	}

	return result
}

// Функция отправки респонса
func sendCommits(writer http.ResponseWriter, request *http.Request) {
	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getCommits(mux.Vars(request)["id"], mux.Vars(request)["date"]))
}

func main() {
	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/{id}", sendCommits).Methods("GET")
	router.HandleFunc("/{id}/", sendCommits).Methods("GET")
	router.HandleFunc("/{id}/{date}", sendCommits).Methods("GET")
	router.HandleFunc("/{id}/{date}/", sendCommits).Methods("GET")

	// Запуск API
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
	//log.Fatal(http.ListenAndServe(":8080", router))
}
