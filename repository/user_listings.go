package repository

import (
	"database/sql"
	"log"

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
	public.***REMOVED***.featured,
	public.***REMOVED***.lat,
	public.***REMOVED***.lon,
	public.***REMOVED***.bathrooms,
	public.***REMOVED***.view_count,
	public.***REMOVED***.street,
	public.***REMOVED***.postal_code,
	public.***REMOVED***.state_province,
	public.***REMOVED***.neighborhood,
	public.***REMOVED***.house_number,
	array_remove(array_agg(public.ad_images.path),NULL) as images 
	FROM public.***REMOVED*** 
	LEFT OUTER JOIN public.ad_images ON (public.***REMOVED***.id = public.ad_images.ad_id) 
	WHERE  public.***REMOVED***.userad_id = $1
	GROUP BY public.***REMOVED***.id
	ORDER BY ***REMOVED***.last_updated ASC`
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
			&ad.Lat,
			&ad.Lon,
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
