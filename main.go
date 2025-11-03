package main

import (
	"fmt"
	"net/http"
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
	fmt.Println("Start")
	http.HandleFunc("/default", handler)

	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	}
}
