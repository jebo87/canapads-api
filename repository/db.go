package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	//"bitbucket.org/jebo87/makako-grpc/ads"
	//"bitbucket.org/jebo/go-postgres-monitor"

	"github.com/lib/pq"
	"gitlab.com/jebo87/makako-grpc/ads"

	//database
	_ "github.com/lib/pq"
)

type ErrorMessage struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

var connInfo string

//InitializeDBConfig used to initialize the configuration for the database
func InitializeDBConfig() string {

	connInfo := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=%v",
		os.Getenv("postgres_host"),
		os.Getenv("postgres_port"),
		os.Getenv("postgres_dbname"),
		os.Getenv("postgres_user"),
		os.Getenv("postgres_password"),
		os.Getenv("postgres_sslmode"))

	return connInfo

}

//GetElasticCount returns the ad count

//GetAdListPB this returns the ads.
//Pagination can be done using offset and limit
func GetAdListPB(db *sql.DB, offset int, limit int) (*ads.AdList, error) {

	adList := &ads.AdList{}

	//if there were errors during database opening we log them here.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := `SELECT public.%[1]v.id, 
	public.%[1]v.title, 
	public.%[1]v.description, 
	public.%[1]v.city, 
	public.%[1]v.country, 
	public.%[1]v.price, 
	public.%[1]v.last_updated, 
	public.%[1]v.rooms, 
	public.%[1]v.property_type, 
	public.%[1]v.userad_id, 
	public.%[1]v.pets, 
	public.%[1]v.furnished, 
	public.%[1]v.garages, 
	public.%[1]v.rent_by_owner, 
	public.%[1]v.published,
	public.%[1]v.last_updated,
	array_remove(array_agg(public.ad_images.path),NULL) as images 
	FROM public.%[1]v 
	LEFT OUTER JOIN public.ad_images ON (public.%[1]v.id = public.ad_images.ad_id) 
	GROUP BY public.%[1]v.id
	ORDER BY %[1]v.last_updated ASC`

	query = fmt.Sprintf(query, os.Getenv("postgres_dbname"))
	if limit == 0 {
		query += " offset " + strconv.Itoa(offset)
	} else if limit != 0 {
		query += " limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset)
	}

	//execute the query and check for errors
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//define an ad to store the ad coming from the database

	//var myTime time.Time
	adList.Ads = []*ads.Ad{}
	//iterate over the results
	for rows.Next() {
		var ad ads.Ad
		//use a pointer to store the title.
		err = rows.Scan(
			&ad.Id,
			&ad.Title,
			&ad.Description,
			&ad.City,
			&ad.Country,
			&ad.Price,
			&ad.PublishedDate,
			&ad.Rooms,
			&ad.PropertyType,
			&ad.UserdadId,
			&ad.Pets,
			&ad.Furnished,
			&ad.Garages,
			&ad.RentByOwner,
			&ad.Published,
			&ad.LastUpdated,
			&ad.Featured,
			&ad.Location,
			&ad.Bathrooms,
			&ad.ViewCount,
			&ad.Street,
			&ad.PostalCode,
			&ad.StateProvince,
			&ad.Neighborhood,
			&ad.HouseNumber,

			(*pq.StringArray)(&ad.Images))
		if err != nil {
			return nil, err
		}

		// ad.PublishedDate = parseDate(ad.PublishedDate, &myTime)

		adList.Ads = append(adList.Ads, &ad)

	}

	return adList, nil

}

//GetAdPB returns the ad matching given ID
func GetAdPB(db *sql.DB, id string) (*ads.Ad, error) {
	if err := db.Ping(); err != nil {
		log.Println("Ping failed", err)
	}

	 
	
	query := `SELECT public.%[1]v.id, 
	public.%[1]v.title, 
	public.%[1]v.description, 
	public.%[1]v.city, 
	public.%[1]v.country, 
	public.%[1]v.price, 
	public.%[1]v.last_updated, 
	public.%[1]v.rooms, 
	public.%[1]v.property_type, 
	public.%[1]v.userad_id, 
	public.%[1]v.pets, 
	public.%[1]v.furnished, 
	public.%[1]v.garages, 
	public.%[1]v.rent_by_owner, 
	public.%[1]v.published,
	public.%[1]v.last_updated,
	public.%[1]v.featured,
	json_build_object('lat',public.%[1]v.lat, 'lon', public.%[1]v.lon) as location, 
	public.%[1]v.bathrooms,
	public.%[1]v.view_count,
	public.%[1]v.street,
	public.%[1]v.postal_code,
	public.%[1]v.state_province,
	public.%[1]v.neighborhood,
	public.%[1]v.house_number,
	array_remove(array_agg(public.ad_images.path),NULL) as images 
	FROM public.%[1]v 
	LEFT OUTER JOIN public.ad_images ON (public.%[1]v.id = public.ad_images.ad_id) 
	WHERE  public.%[1]v.id = $1
	GROUP BY public.%[1]v.id
	ORDER BY %[1]v.last_updated ASC`
	query = fmt.Sprintf(query, os.Getenv("postgres_dbname"))

	row := db.QueryRow(query, id)

	log.Println("Fetching ad", id, " from database")

	//var myTime time.Time
	var ad Ad
	err := row.Scan(&ad.ID,
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
		&ad.RentByOwner,
		&ad.Published,
		&ad.LastUpdated,
		&ad.Featured,
		&ad.Location,
		&ad.Bathrooms,
		&ad.ViewCount,
		&ad.Street,
		&ad.PostalCode,
		&ad.StateProvince,
		&ad.Neighborhood,
		&ad.HouseNumber,
		(*pq.StringArray)(&ad.Images))

		

	if err != nil {
		log.Println("error while trying to scan",err)
		return &ads.Ad{}, errors.New("Ad not found")
	}
	

	//ad.PublishedDate = parseDate(ad.PublishedDate, &myTime)
	log.Println("returning ad ", id, " ", ad.Title, " ", ad.PublishedDate)
	if ad.ID == 0 {
		return &ads.Ad{}, errors.New("Ad not found")
	}

	return ToProto(ad, &ads.Ad{}), err

}

//InsertAd Inserts a new add into the database
func InsertAd(ad Ad) {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", connInfo)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	checkErr(err, "fatal")
	var id int

	//prepare the query statement
	query := `INSERT INTO %[1]v (
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

	query = fmt.Sprintf(query, os.Getenv("postgres_dbname"))

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
	db, err := sql.Open("postgres", connInfo)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	checkErr(err, "fatal")

	//prepare the query string
	query := "DELETE FROM %[1]v where id = $1"

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
