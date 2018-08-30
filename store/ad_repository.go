package store

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"bitbucket.org/jebo87/makako-grpc/ads"

	"github.com/lib/pq"

	//database
	_ "github.com/lib/pq"
	yaml "gopkg.in/yaml.v2"
)

//Config struct to handle the YAML configuration
type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	}
	Elastic struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	DBInfo struct {
		TargetTable  string `yaml:"table"`
		Schema       string `yaml:"schema"`
		MonitorTable string `yaml:"monitor"`
	}
}

var conf Config
var connInfo string

//InitializeDB used to initialize the configuration for the database
func InitializeDB() {
	//connection string for the database
	conf = loadConfig()
	connInfo = "host=" + conf.Postgres.Host +
		" dbname=" + conf.Postgres.DBName +
		" user=" + conf.Postgres.User +
		" password=" + conf.Postgres.Password +
		" sslmode=" + conf.Postgres.SSLMode
}

//loadConfig loads the configuration from a yaml file
func loadConfig() (conf Config) {
	configFile, err := ioutil.ReadFile("config/conf.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}
	return

}

//parseDate parses the date coming from the database to a format
//compatible with protobuffers
func parseDate(pdate *ads.Date, fecha time.Time) *ads.Date {
	pdate = &ads.Date{}
	pdate.Year = int32(fecha.Year())

	pdate.Month = int32(fecha.Month())
	pdate.Day = int32(fecha.Day())

	return pdate

}

//GetAdListPB this returns the ads.
//Pagination can be done using offset and limit
func GetAdListPB(offset int, limit int) (*ads.AdList, error) {

	adList := &ads.AdList{}

	//open the connection and store errors in err
	db, err := sql.Open("postgres", connInfo)

	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := `SELECT public.***REMOVED***.id, 
	public.***REMOVED***.title, 
	public.***REMOVED***.description, 
	public.***REMOVED***.city, 
	public.***REMOVED***.country, 
	public.***REMOVED***.price, 
	public.***REMOVED***.last_updated, 
	public.***REMOVED***.rooms, 
	public.***REMOVED***.property_type, 
	public.***REMOVED***.userad_id, 
	public.***REMOVED***.pets, 
	public.***REMOVED***.furnished, 
	public.***REMOVED***.garages, 
	public.***REMOVED***.rent_by_owner, 
	public.***REMOVED***.published,
	public.***REMOVED***.last_updated,
	array_agg(public.ad_images.path) as images 
	FROM public.***REMOVED*** 
	LEFT OUTER JOIN public.ad_images ON (public.***REMOVED***.id = public.ad_images.ad_id) 
	GROUP BY public.***REMOVED***.id
	ORDER BY ***REMOVED***.last_updated ASC`

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

	//define an ad to store the ad coming from the database

	var myTime time.Time
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
			&myTime,
			&ad.Rooms,
			&ad.PropertyType,
			&ad.UserdadId,
			&ad.Pets,
			&ad.Furnished,
			&ad.Garages,
			&ad.RentByOwner,
			&ad.Published,
			&ad.LastUpdated,

			(*pq.StringArray)(&ad.Images))
		if err != nil {
			return nil, err
		}
		ad.PublishedDate = parseDate(ad.PublishedDate, myTime)

		adList.Ads = append(adList.Ads, &ad)

	}

	return adList, nil

}

//GetAd returns the ad matching given ID
func GetAdPB(id string) (*ads.Ad, error) {
	//open the connection and store errors in err
	db, err := sql.Open("postgres", connInfo)
	//defer the database close
	defer db.Close()

	//if there were errors during database opening we log them here.
	if err != nil {
		log.Fatal(err)
	}

	//modify the query depending on the number of ads to display
	query := `SELECT public.***REMOVED***.id, 
	public.***REMOVED***.title, 
	public.***REMOVED***.description, 
	public.***REMOVED***.city, 
	public.***REMOVED***.country, 
	public.***REMOVED***.price, 
	public.***REMOVED***.last_updated, 
	public.***REMOVED***.rooms, 
	public.***REMOVED***.property_type, 
	public.***REMOVED***.userad_id, 
	public.***REMOVED***.pets, 
	public.***REMOVED***.furnished, 
	public.***REMOVED***.garages, 
	public.***REMOVED***.rent_by_owner, 
	public.***REMOVED***.published,
	public.***REMOVED***.last_updated,
	array_agg(public.ad_images.path) as images 
	FROM public.***REMOVED*** 
	LEFT OUTER JOIN public.ad_images ON (public.***REMOVED***.id = public.ad_images.ad_id) 
	WHERE  public.***REMOVED***.id = $1
	GROUP BY public.***REMOVED***.id
	ORDER BY ***REMOVED***.last_updated ASC`

	row := db.QueryRow(query, id)
	var myTime time.Time
	var ad ads.Ad
	row.Scan(&ad.Id,
		&ad.Title,
		&ad.Description,
		&ad.City,
		&ad.Country,
		&ad.Price,
		&myTime,
		&ad.Rooms,
		&ad.PropertyType,
		&ad.UserdadId,
		&ad.Pets,
		&ad.Furnished,
		&ad.Garages,
		&ad.RentByOwner,
		&ad.Published,
		&ad.LastUpdated,

		(*pq.StringArray)(&ad.Images))

	checkErr(err, "panic")

	ad.PublishedDate = parseDate(ad.PublishedDate, myTime)

	if ad.Id == 0 {

		return &ad, errors.New("Ad not found")
	}

	return &ad, err

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
	db, err := sql.Open("postgres", connInfo)
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
