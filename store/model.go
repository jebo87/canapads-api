package store

import (
	"strconv"
	"time"
)

// Ad basic structure
type Ad struct {
	ID            int
	Title         string
	Description   string
	City          string
	Country       string
	Price         int
	PublishedDate time.Time
	Rooms         int
	PropertyType  string
	UserAdID      int
	Pets          bool
	Furnished     bool
	Garages       int
	RentByOwner   bool
}

//AdToString returns a string for the ad
func AdToString(ad Ad) string {
	return strconv.Itoa(ad.ID) + " " + ad.Title
}
