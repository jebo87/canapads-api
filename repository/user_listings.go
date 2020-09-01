package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var list ads.AdList

//GetUserListings returns all listings associated to a given user
func GetUserListings(db *sql.DB, userID string) (adList *ads.AdList, err error) {
	log.Println("GetUserListings for ", userID)
	//if there were errors during database opening we log them here.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
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
	public.%[1]v.lat,
	public.%[1]v.lon,
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
	WHERE  public.%[1]v.userad_id = $1
	GROUP BY public.%[1]v.id
	ORDER BY %[1]v.last_updated ASC`
	query = fmt.Sprintf(query, os.Getenv("postgres_dbname"))

	log.Println(query)
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var tempList []*ads.Ad
	for rows.Next() {
		var ad ads.Ad
		err = rows.Scan(&ad.Id,
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
			log.Println(err)

			return nil, err
		}

		tempList = append(tempList, &ad)
		log.Println(ad.Id)
	}
	list.Ads = tempList
	return &list, nil
}
