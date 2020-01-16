package store

import (
	"strconv"

	"gitlab.com/jebo87/makako-grpc/ads"
)

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
	UserAdID      int      `json:"userad_id"`
	Pets          int      `json:"pets"`
	Furnished     bool     `json:"furnished"`
	Garages       int      `json:"garages"`
	RentByOwner   bool     `json:"rent_by_owner"`
	Published     bool     `json:"published"`
	LastUpdated   string   `json:"last_updated"`
	Featured      int      `json:"featured"`
	Lat           float64  `json:"lat"`
	Lon           float64  `json:"lon"`
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
	pb.UserdadId = int32(ad.UserAdID)
	pb.Pets = int32(ad.Pets)
	pb.Furnished = ad.Furnished
	pb.Garages = int32(ad.Garages)
	pb.RentByOwner = ad.RentByOwner
	pb.Published = ad.Published
	pb.Featured = int32(ad.Featured)
	pb.Lat = float64(ad.Lat)
	pb.Lon = float64(ad.Lon)
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
