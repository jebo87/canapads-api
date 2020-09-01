package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gitlab.com/jebo87/makako-grpc/ads"
)

// Location basic structure
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (s Location) Value() (driver.Value, error) {

	return fmt.Sprintf(`{lat:"%v",lon:"%v"}`, s.Lat, s.Lon), nil
}

func (s *Location) Scan(src interface{}) (err error) {
	var location Location
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &location)
	case []byte:
		err = json.Unmarshal(src.([]byte), &location)
	default:
		return errors.New("Incompatible type for location")
	}
	if err != nil {
		return
	}
	*s = location
	return nil
}

// Ad basic structure
type Ad struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	City          string   `json:"city"`
	Country       string   `json:"country"`
	Price         int      `json:"price"`
	PublishedDate string   `json:"published_date"`
	Rooms         int      `json:"rooms"`
	PropertyType  string   `json:"property_type"`
	UserAdID      string   `json:"userad_id"`
	Pets          int      `json:"pets"`
	Furnished     bool     `json:"furnished"`
	Garages       int      `json:"garages"`
	RentByOwner   bool     `json:"rent_by_owner"`
	Published     bool     `json:"published"`
	LastUpdated   string   `json:"last_updated"`
	Featured      int      `json:"featured"`
	Location      Location `json:"location"`
	Bathrooms     int      `json:"bathrooms"`
	ViewCount     int      `json:"view_count"`
	Street        string   `json:"street"`
	PostalCode    string   `json:"postal_code"`
	StateProvince string   `json:"state_province"`
	Neighborhood  string   `json:"neighborhood"`
	HouseNumber   string   `json:"house_number"`
	Gym           bool     `json:"gym"`
	Pool          bool     `json:"pool"`
	Images        []string `json:"images"`
}

type Filter struct {
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	City              string  `json:"city"`
	Country           string  `json:"country"`
	PriceLow          int     `json:"price_low"`
	PriceHigh         int     `json:"price_high"`
	PublishedDateLow  string  `json:"published_date_low"`
	PublishedDateHigh string  `json:"published_date_high"`
	RoomsLow          int     `json:"rooms_low"`
	RoomsHigh         int     `json:"rooms_high"`
	PropertyType      string  `json:"property_type"`
	Pets              int     `json:"pets"`
	Furnished         bool    `json:"furnished"`
	Garages           int     `json:"garages"`
	RentByOwner       bool    `json:"rent_by_owner"`
	Lat               string  `json:"lat"`
	Lon               string  `json:"lon"`
	Bathrooms         int     `json:"bathrooms"`
	PostalCode        string  `json:"postal_code"`
	StateProvince     string  `json:"state_province"`
	Neighborhood      string  `json:"neighborhood"`
	Gym               bool    `json:"gym"`
	Pool              bool    `json:"pool"`
	HasImages         bool    `json:"hasImages"`
	From              int     `json:"from"`
	Size              int     `json:"size"`
	SearchParam       string  `json:"searchParam"`
	Points            Polygon `json:"points"`
}
type Polygon struct {
	points []Location
}

//AdToString returns a string for the ad
func AdToString(ad Ad) string {
	return strconv.Itoa(ad.ID) + " " + ad.Title
}

//ToProto maps an ad to a protobuf ad
func ToProto(ad Ad, pb *ads.Ad) *ads.Ad {
	pb.Id = int32(ad.ID)
	pb.Title = ad.Title
	pb.Description = ad.Description
	pb.City = ad.City
	pb.Country = ad.Country
	pb.Price = int32(ad.Price)
	pb.PropertyType = ad.PropertyType
	pb.PublishedDate = ad.PublishedDate
	pb.Rooms = int32(ad.Rooms)
	pb.UserdadId = ad.UserAdID
	pb.Pets = int32(ad.Pets)
	pb.Furnished = ad.Furnished
	pb.Garages = int32(ad.Garages)
	pb.RentByOwner = ad.RentByOwner
	pb.Published = ad.Published
	pb.Featured = int32(ad.Featured)
	pb.Location = &ads.Location{Lat: ad.Location.Lat, Lon: ad.Location.Lon}
	pb.Bathrooms = int32(ad.Bathrooms)
	pb.ViewCount = int32(ad.ViewCount)
	pb.Street = ad.Street
	pb.PostalCode = ad.PostalCode
	pb.Neighborhood = ad.Neighborhood
	pb.HouseNumber = ad.HouseNumber
	pb.Gym = ad.Gym
	pb.Pool = ad.Pool
	pb.Images = ad.Images
	//log.Println(pb)
	//log.Println("**************************************")
	//log.Println(ad)
	return pb
}

type elasticSearchAd struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"succesful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int     `json:"total"`
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float32 `json:"_score"`
			Source Ad      `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
type elasticSearchAdCount struct {
	Count  int `json:"count"`
	Shards struct {
		Total      int `json:"total"`
		Successful int `json:"succesful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}
