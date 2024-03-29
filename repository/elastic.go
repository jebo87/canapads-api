package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.com/jebo87/makako-gateway/httputils"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var remoteAddres string

//GetElasticCount returns the total quantity of ads
func GetElasticCount() (*ads.AdCount, error) {
	var resp *http.Response
	var err error

	resp, err = http.Get("http://" + os.Getenv("elastic.host") + ":" + os.Getenv("elastic.port") + "/ads/ad/_count")

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	results := &elasticSearchAdCount{}
	log.Printf("[%v] Reading response from Elastic", remoteAddres)
	body, err := ioutil.ReadAll(resp.Body)
	count := &ads.AdCount{}
	if resp.StatusCode == 200 {

		err = json.Unmarshal(body, &results)
		log.Println(string(body))
		if err != nil {
			log.Println("ERROR", err)
		}
		count.Count = int32(results.Count)
	} else {
		err = errors.New(string(body))
	}

	return count, err

}
func prepareQueryParam(searchTerm string) string {
	queryParam := ""
	if searchTerm != "" {
		queryParam = `
"should": [
	{
		"match": {
			"description": {
				"query": %q
			
			}
		}
	},{
		"match": {
			"title": {
				"query": %q,
				"boost":2
			
			}
		}
	}
],"minimum_should_match" : 1,
`

		queryParam = fmt.Sprintf(queryParam, searchTerm, searchTerm)
	} else {
		queryParam = ""
	}
	return queryParam
}

//prepareBody Preparse the json body to be submited to elastic
//the query is prepared depending on the filters received from the gateway
func prepareBody(searchTerm string, singleValuefilters map[string]string, fromSize map[string]string, priceRange string, polygonFilter string) string {

	queryParam := prepareQueryParam(searchTerm)

	search := `
{
"_source": [
		"id",
		"title",
		"description",
		"published_date",
		"price",
		"last_updated",
		"images",
		"location"
	],
	"from": %v,
	"size": %v,
	"query": {
		"bool": {
			%v
			"filter": [
				%v
			]
		}
	}
}
`

	var terms []string

	for k, v := range singleValuefilters {
		terms = append(terms, fmt.Sprintf(`{"term": {"%v": %v}}`, k, v))
	}

	//this validation is needed to parsing issues when the polygon filter is empty
	if polygonFilter != "" {
		terms = append(terms, priceRange, polygonFilter)
	} else {
		terms = append(terms, priceRange)

	}
	//put everything into the json to be sent.
	body := fmt.Sprintf(search, fromSize["from"], fromSize["size"], queryParam, strings.Join(terms, ","))
	log.Printf("[%v] Query object sent to elastic:", remoteAddres)
	//print the json in the console for troubleshooting purposes
	log.Println("\n", httputils.JSONPrettyPrint(body))
	return body
}

func preparePolygonFilter(filter *ads.Filter) string {
	coordinates := ""
	if len(filter.GetPolygon().GetPoints()) > 0 {

		for i, v := range filter.GetPolygon().GetPoints() {

			coordinates += fmt.Sprintf(`
			{
				"lon":%v,
				"lat":%v
			}`, v.Lon, v.Lat)
			if i < len(filter.GetPolygon().GetPoints())-1 {
				coordinates += ","
			}
		}

	} else {
		return ""
	}
	polygonFilter := `
	{
		"geo_polygon": {
			"location": {
				"points": [
					%v
				]
			}
		}
	}
	`

	return fmt.Sprintf(polygonFilter, coordinates)
}

func preparePriceRangeFilter(filter *ads.Filter) string {
	price := make(map[string]int)
	//default values for price search
	//TODO: this needs to be revised if we ever decide to list houses for sale
	price["gte"] = 0
	price["lte"] = 1000000
	if filter.GetPriceLow() != nil {
		price["gte"] = int(filter.GetPriceLow().GetValue())
	}
	if filter.GetPriceHigh() != nil {
		price["lte"] = int(filter.GetPriceHigh().GetValue())
	}

	priceRange := `{
					"range": 
						{"price": 
							{"gte": %v, "lte": %v} 
						} 
					}`
	priceRange = fmt.Sprintf(priceRange, price["gte"], price["lte"])

	return priceRange

}

func prepareFromSizeFilter(filter *ads.Filter) map[string]string {
	myFromSizeMap := make(map[string]string)
	if filter.GetSize() != nil {
		myFromSizeMap["size"] = strconv.Itoa(int(filter.GetSize().GetValue()))
	}
	if filter.GetFrom() != nil {
		myFromSizeMap["from"] = strconv.Itoa(int(filter.GetFrom().GetValue()))
	}

	return myFromSizeMap
}

func prepareSingleValueFilters(filter *ads.Filter) map[string]string {
	log.Printf("[%v] preparing single value filters", remoteAddres)

	myFilterMap := make(map[string]string)
	if filter.GetGym() != nil {
		myFilterMap["gym"] = fmt.Sprintf("%v", filter.GetGym().GetValue())
	}
	if filter.GetPets() != nil {
		myFilterMap["pets"] = fmt.Sprintf("%v", filter.GetPets().GetValue())
	}
	if filter.GetPool() != nil {
		myFilterMap["pool"] = fmt.Sprintf("%v", filter.GetPool().GetValue())
	}
	if filter.GetCity() != nil {
		myFilterMap["city"] = fmt.Sprintf("%v", filter.GetCity().GetValue())
	}
	if filter.GetCountry() != nil {
		myFilterMap["country"] = fmt.Sprintf("%v", filter.GetCountry().GetValue())
	}
	if filter.GetPropertyType() != nil {
		myFilterMap["property_type"] = fmt.Sprintf("%q", filter.GetPropertyType().GetValue())
	}
	if filter.GetFurnished() != nil {
		myFilterMap["furnished"] = fmt.Sprintf("%v", filter.GetFurnished().GetValue())
	}
	if filter.GetRentByOwner() != nil {
		myFilterMap["rent_by_owner"] = fmt.Sprintf("%v", filter.GetRentByOwner().GetValue())
	}
	if filter.GetStateProvince() != nil {
		myFilterMap["province"] = fmt.Sprintf("%v", filter.GetStateProvince().GetValue())
	}
	if filter.GetNeighborhood() != nil {
		myFilterMap["neighborhood"] = fmt.Sprintf("%v", filter.GetNeighborhood().GetValue())
	}
	if filter.GetGarages() != nil {
		myFilterMap["garages"] = fmt.Sprintf("%v", filter.GetGarages().GetValue())
	}
	if filter.GetRooms() != nil {
		myFilterMap["rooms"] = fmt.Sprintf("%v", filter.GetRooms().GetValue())
	}

	if filter.GetBathrooms() != nil {
		myFilterMap["bathrooms"] = fmt.Sprintf("%v", filter.GetBathrooms().GetValue())
	}

	// google.protobuf.StringValue lat = 16;
	// google.protobuf.StringValue lon = 17;
	// google.protobuf.Int32Value bathrooms = 18;

	// google.protobuf.BoolValue hasImages = 24;
	// google.protobuf.StringValue published_date_low = 7;
	// google.protobuf.StringValue published_date_high = 8;
	// google.protobuf.Int32Value rooms_low = 9;
	// google.protobuf.Int32Value rooms_high= 10;
	//log.Printf("filter %v", filter)

	return myFilterMap
}

//SearchElastic serach in elastic search
func SearchElastic(filter *ads.Filter, remoteAddr string) (*ads.SearchResponse, error) {
	//log.Println(filter)
	remoteAddres = remoteAddr
	//this map will contain all the applicable filters received in the request
	//we must validate each type of filter to be able to set them properly for elasticSearch
	myFilterMap := prepareSingleValueFilters(filter)
	fromSize := prepareFromSizeFilter(filter)
	priceRange := preparePriceRangeFilter(filter)
	polygonFilter := preparePolygonFilter(filter)

	requestBody := []byte(prepareBody(filter.GetSearchParam().GetValue(), myFilterMap, fromSize, priceRange, polygonFilter))
	adList := &ads.AdList{}
	searchResponse := &ads.SearchResponse{}
	var err error
	var req *http.Request

	req, _ = http.NewRequest("POST", "http://"+os.Getenv("elastic_host")+":"+os.Getenv("elastic_port")+"/_search", bytes.NewBuffer(requestBody))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%v] Error posting request %v", remoteAddres, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		log.Printf("[%v] Connected to Elastic...", remoteAddres)
		//create a struct to hold the values from the response
		results := &elasticSearchAd{}

		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("[%v] Reading response from Elastic...", remoteAddres)

		err = json.Unmarshal(body, &results)
		log.Printf("[%v] took :%v, timed_out:%v ,hits: %v", remoteAddres, results.Took, results.TimedOut, results.Hits.Total)
		adList.Ads = []*ads.Ad{}
		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Printf("[%v] Translating ads to protobuf...", remoteAddres)

		//convert the ads to protobuf and add them to the adList that will be returned
		for _, ad := range results.Hits.Hits {
			adPB := &ads.Ad{}
			adPB = ToProto(ad.Source, adPB)
			adList.Ads = append(adList.Ads, adPB)

		}
		log.Printf("[%v] Listings loaded", remoteAddres)

		searchResponse.List = adList
		searchResponse.Count = int32(results.Hits.Total)

		return searchResponse, nil
	}

	//TODO: return custom errors like the ones coming from elastic, this will help troubleshoot in case of problems
	// {
	// 	"error": "Incorrect HTTP method for uri [/ads/ad/?sort] and method [GET], allowed: [POST]",
	// 	"status": 405
	// 	}
	//log.Println(resp)
	return searchResponse, errors.New("status" + strconv.Itoa(resp.StatusCode) + " MakakoLabs: There was a problem while procesing your request")

}
