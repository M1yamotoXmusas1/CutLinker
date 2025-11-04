package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// 1. Получить запрос
// 2. Обработать запрос
// 3. сгенерить ссылку
// 4. занести в базу данных
// 5. создать страничку
// 6. profit

func handler(w http.ResponseWriter, r *http.Request) {
	str := "Hello World!"
	b := []byte(str)

	_, err := w.Write(b)
	if err != nil {
		fmt.Println("Во время записи HTTP ответа произошла ошибка:", err.Error())
	} else {
		fmt.Println("Корректно обработан HTTP запрос!")
	}
}

func main() {
	connStr := "postgres://postgres:AZRAELBEATS@localhost:5431/golinkdb?sslmode=disable"

	db, errdb := sql.Open("postgres", connStr)

	if errdb != nil {
		log.Fatal(errdb)
	}
	if errdb = db.Ping(); errdb != nil {
		log.Fatal("Ошибка Open db:", errdb)
	}
	fmt.Println("Дб запустилась!")
	defer db.Close()

	http.HandleFunc("/createlink", handler)
	http.HandleFunc("/redirect", handler)

	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	}
}
