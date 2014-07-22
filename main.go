package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"net/http"
)

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", "dbname=lesson4 sslmode=disable")
	panicIf(err)
	return db
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	m := martini.Classic()
	//m.Get("/", func() string {
	//	//return "Hello world! AHAHAHAHA"
	//	return `
	//	{ message: "This is IT" }
	//	`
	//})
	m.Map(setupDB())
	m.Use(render.Renderer())
	m.Get("/", func(req *http.Request, db *sql.DB, rw http.ResponseWriter, r render.Render) {
		search := "%" + req.URL.Query().Get("search") + "%"
		rows, err := db.Query(`SELECT title, author, description 
                           FROM books 
                           WHERE title ILIKE $1
                           OR author ILIKE $1
                           OR description ILIKE $1`, search)
		panicIf(err)
		defer rows.Close()

		var title, author, description string
		var results []map[string]interface{}

		for rows.Next() {
			err := rows.Scan(&title, &author, &description)
			panicIf(err)

			results = append(results, map[string]interface{}{
				"Title":       title,
				"Author":      author,
				"Description": description,
			})
		}

		r.JSON(200, results)
	})
	m.Use(martini.Static("stat"))
	m.Run()
}
