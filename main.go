package main

import (
	"fmt"
	"time"

	"bitbucket.org/jebo87/makako-api/store"
)

func main() {
	fmt.Println("Loading makako API server...")
	store.Initialize()
	// store.GetAdTitles(0, 0)
	// testAddAd()
	store.DeleteAd(9)

}

func testAddAd() {
	var ad store.Ad
	ad.Title = "prueba 1"
	ad.Description = "Description 1"
	ad.City = "Montreal"
	ad.Country = "Canada"
	ad.Price = 660
	ad.PublishedDate = time.Now()
	ad.PropertyType = "apartment"
	ad.Rooms = 4
	ad.UserAdID = 1

	store.InsertAd(ad)

}
