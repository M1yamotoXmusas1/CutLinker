package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// 1. Получить запрос
// 2. Обработать запрос
// 3. сгенерить ссылку
// 4. занести в базу данных
// 5. создать страничку
// 6. profit

type shortlink struct {
	baseurl  string
	shorturl string
	clicks   int
}

func createLinkTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS link (
		id SERIAL PRIMARY KEY,
		baseurl VARCHAR(100) NOT NULL,
		shorturl VARCHAR(100) NOT NULL,
		clicks NUMERIC(6),
		created timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatal("Fail to create a table: ", err)
	}
}

func insertLinkTable(db *sql.DB, link shortlink) int {
	query := `INSERT INTO link (baseurl, shorturl, clicks)
		VALUES ($1, $2, $3) RETURNING id`

	var pk int
	err := db.QueryRow(query, link.baseurl, link.shorturl, link.clicks).Scan(&pk)
	if err != nil {
		log.Fatal("Fail to insert into table", err)
	}
	fmt.Printf("pk: %v\n", pk)
	return pk
}

func generateRandomString(lenght int) string {
	// 16 символов в байт -> 12 в base64
	lenght = (lenght * 3) / 4
	bytes := make([]byte, lenght)
	rand.Read(bytes)

	return base64.URLEncoding.EncodeToString(bytes)
}

func CreateURLHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpRequestBody, err_request := io.ReadAll(r.Body)

		if err_request != nil {
			fmt.Println("Fail to read HTTP body: ", err_request)
		}

		url := string(httpRequestBody)

		if !strings.HasPrefix(url, "https://") {
			fmt.Printf("URL is not valid %v\n", url)
			return
		}

		shortUrl := "https://localhost:11111/" + generateRandomString(8)

		link := shortlink{url, shortUrl, 0}
		insertLinkTable(db, link)

		response := []byte(shortUrl)
		_, err_write := w.Write(response)
		if err_write != nil {
			fmt.Println("Fail to write HTTP response: ", err_write)
		}
	}
}

func main() {
	connStr := "postgres://postgres:AZRAELBEATS@localhost:5431/golinkdb?sslmode=disable"

	db, err_db := sql.Open("postgres", connStr)

	if err_db != nil {
		log.Fatal("Fail to open database: ", err_db)
	}
	if err_db = db.Ping(); err_db != nil {
		log.Fatal("Fail ping() database: ", err_db)
	}
	fmt.Println("Дб запустилась!")
	defer db.Close()
	createLinkTable(db)

	http.HandleFunc("/createURL", CreateURLHandler(db))
	//http.HandleFunc("/redirect", CreateURLHandler)

	err := http.ListenAndServe(":11111", nil)
	if err != nil {
		fmt.Println("Произошла ошибка:", err)
	}
}
