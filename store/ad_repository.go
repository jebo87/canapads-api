package store

import (
	"database/sql"
	"encoding/json"
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

//InitializeDB used to initialize the configuration for the database
func InitializeDB() {
	//connection string for the database
	conf.connStr = "postgres://postgres:test123@dbads.makakolabs.ca/***REMOVED***?sslmode=disable"
}

//GetAdTitles this returns the ads.
//Pagination can be done using offset and limit
func GetAdTitles(offset int, limit int) []byte {

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
	checkErr(err, "panic")

	//define an ad to store the ad coming from the database
	var ad Ad

	ads := make([]Ad, 0)

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
			&ad.UserAdID,
			&ad.Pets,
			&ad.Furnished,
			&ad.Garages,
			&ad.RentByOwner)
		checkErr(err, "panic")

		//append the ad to the slice
		ads = append(ads, ad)

	}
	// convert the slice into a byte array and return it
	data, _ := json.Marshal(ads)
	return data

}

//GetAd returns the ad matching given ID
func GetAd(id string) []byte {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", conf.connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := "SELECT * FROM ***REMOVED*** WHERE ID=$1"

	row := db.QueryRow(query, id)

	var ad Ad
	row.Scan(&ad.ID,
		&ad.Title,
		&ad.Description,
		&ad.City,
		&ad.Country,
		&ad.Price,
		&ad.PublishedDate,
		&ad.Rooms,
		&ad.PropertyType,
		&ad.UserAdID,
		&ad.Pets,
		&ad.Furnished,
		&ad.Garages,
		&ad.RentByOwner)

	if ad.ID == 0 {
		return nil
	}

	// convert the ad into a byte array and return it
	data, _ := json.Marshal(ad)
	return data

}

//InsertAd Inserts a new add into the database
func InsertAd(ad Ad) {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", conf.connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	checkErr(err, "fatal")
	var id int

	//prepare the query statement
	query := `INSERT INTO ***REMOVED*** (
		title,
		description,
		city,
		country,
		price,
		published_date,
		rooms,
		property_type,
		userad_id,
		pets,	
		furnished,	  
		garages,		  
		rent_by_owner   
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`

	//execute the statement adding the values from the ad and return the id of the newly created ad
	err = db.QueryRow(
		query,
		ad.Title,
		ad.Description,
		ad.City,
		ad.Country,
		ad.Price,
		ad.PublishedDate,
		ad.Rooms,
		ad.PropertyType,
		ad.UserAdID,
		ad.Pets,
		ad.Furnished,
		ad.Garages,
		ad.RentByOwner).Scan(&id)
	checkErr(err, "panic")
	fmt.Println(id)

}

//DeleteAd deletes an ad from the database
func DeleteAd(index int) {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", conf.connStr)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	checkErr(err, "fatal")

	//prepare the query string
	query := "DELETE FROM ***REMOVED*** where id = $1"

	//execute the query
	_, err = db.Exec(query, index)
	checkErr(err, "panic")

}

func checkErr(err error, errorType string) {
	if err != nil {
		if errorType == "panic" {
			panic(err)
		} else if errorType == "fatal" {
			log.Fatal(err)
		}
	}
}
