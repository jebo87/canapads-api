package store

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	//database
	_ "github.com/lib/pq"
)

//Config structure
type Config struct {
	connStr string
}

var conf Config

//Initialize used to initialize the configuration for the database
func Initialize() {
	//connection string for the database
	conf.connStr = "postgres://postgres:test123@dbads.makakolabs.ca/***REMOVED***?sslmode=disable"
}

//GetAdTitles this returns all ads
func GetAdTitles(offset int, limit int) {

	//open the connection and store errors in err
	db, err := sql.Open("postgres", conf.connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := "SELECT * FROM ***REMOVED***"
	if limit == 0 {
		query += " offset " + strconv.Itoa(offset)
	} else if limit != 0 {
		query += " limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset)
	}

	//execute the query and check for errors
	rows, err := db.Query(query)
	checkErr(err)

	//define an ad to store the ad coming from the database
	var ad Ad

	//iterate over the results
	for rows.Next() {

		//use a pointer to store the title.
		err = rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Description,
			&ad.City,
			&ad.Country,
			&ad.Price,
			&ad.PublishedDate,
			&ad.Rooms,
			&ad.PropertyType,
			&ad.UserAdID)
		checkErr(err)

		fmt.Println(AdToString(ad))
	}

}

//InsertAd Inserts a new add into the database
func InsertAd(ad Ad) {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", conf.connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	var id int

	//prepare the query statement
	query := "INSERT INTO ***REMOVED*** (title,description,city,country,price,publishedDate,rooms,property_type,useradid) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"

	//execute the statement adding the values from the ad and return the id of the newly created ad
	err = db.QueryRow(query, ad.Title, ad.Description, ad.City, ad.Country, ad.Price, ad.PublishedDate, ad.Rooms, ad.PropertyType, ad.UserAdID).Scan(&id)
	checkErr(err)
	fmt.Println(id)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func hola() {}
