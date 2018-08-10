package store

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	//database
	_ "github.com/lib/pq"
)

//GetAdTitles this returns all ads
func GetAdTitles(offset int, limit int) {

	//connection string for the database
	connStr := "postgres://postgres:test123@dbads.makakolabs.ca/***REMOVED***?sslmode=disable"

	//open the connection and store errors in err
	db, err := sql.Open("postgres", connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := "SELECT id, title, description FROM ***REMOVED***"
	if limit == 0 {
		query += " offset " + strconv.Itoa(offset)
	} else if limit != 0 {
		query += " limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset)
	}

	//execute the query and check for errors
	rows, err := db.Query(query)
	checkErr(err)

	var id int
	var title string
	var description string
	//iterate over the results
	for rows.Next() {

		//use a pointer to store the title.
		err = rows.Scan(&id, &title, &description)
		checkErr(err)

		fmt.Println(id, title, description)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func hola() {}
