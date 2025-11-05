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

// 1. added error handler to func getRandomString
// 2. error print with fmt replaced with log

const IP = "http://localhost:11111"

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
	return pk
}

func generateRandomString(length int) string {
	// 16 символов в байт -> 12 в base64
	length = (length * 3) / 4
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)

	if err != nil {
		log.Println("Failed to generate string: ", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func CreateURLHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpRequestBody, err_request := io.ReadAll(r.Body)

		if err_request != nil {
			log.Println("Fail to read HTTP body: ", err_request)
		}

		url := string(httpRequestBody)

		if !strings.HasPrefix(url, "https://") {
			log.Println("URL is not valid: ", url)
			return
		}

		shortUrl := "http://localhost:11111/" + generateRandomString(8)

		link := shortlink{url, shortUrl, 0}
		insertLinkTable(db, link)

		response := []byte(shortUrl)
		_, err_write := w.Write(response)
		if err_write != nil {
			log.Println("Fail to write HTTP response: ", err_write)
		}
	}
}

func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			return
		}
		query := `SELECT baseurl FROM link WHERE shorturl = $1`
		baseurl := ""
		shorturl := IP + r.URL.Path
		err := db.QueryRow(query, shorturl).Scan(&baseurl)
		fmt.Println("CHLEEEEN", shorturl)
		fmt.Println("ADSP[aD]", baseurl)
		if err != nil {
			log.Println("Failed to find url in database: ", err)
		}
		log.Println("Alright! baseurl: ", baseurl)
		http.Redirect(w, r, baseurl, http.StatusSeeOther)
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
	http.HandleFunc("/", RedirectHandler(db))

	err := http.ListenAndServe(":11111", nil)
	if err != nil {
		log.Println("Произошла ошибка:", err)
	}
}
