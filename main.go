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
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения информации с сайта
func getInfo(username string, date string) User {
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

	// Так выглядит html одной ячейки:
	// <rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1">

	// Весь html страницы в формате string
	pageStr := string(body[140000:225000])

	// Указатель на ячейку нужной даты
	i := strings.Index(pageStr, "data-date=\""+date)

	// Проверка на существование нужной ячейки
	if i != -1 {
		for ; pageStr[i] != '<'; i-- {
			// Доводит i до начала кода ячейки
		}

		// Получение параметров ячейки
		values := strings.FieldsFunc(pageStr[i:i+155], func(r rune) bool {
			return r == '"'
		})

		// Запись и обработка нужной информации
		dataLevel, _ := strconv.Atoi(values[19])
		commits, _ := strconv.Atoi(values[15])

		// Возвращение обработанной информации
		return User{
			Date:     date,
			Username: username,
			Commits:  commits,
			Color:    dataLevel,
		}
	}

	// Если пользователь не найден, API вернёт пустую структуру
	return User{
		Date:     date,
		Username: username,
		Commits:  0,
		Color:    0,
	}
}

// Функция отправки респонса
func sendInfo(writer http.ResponseWriter, request *http.Request) {
	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getInfo(mux.Vars(request)["id"], mux.Vars(request)["date"]))
}

func main() {
	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/{id}", sendInfo).Methods("GET")
	router.HandleFunc("/{id}/", sendInfo).Methods("GET")
	router.HandleFunc("/{id}/{date}", sendInfo).Methods("GET")
	router.HandleFunc("/{id}/{date}/", sendInfo).Methods("GET")

	// Запуск API
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
	// log.Fatal(http.ListenAndServe(":8080", router))
}
