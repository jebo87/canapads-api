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
func GetElasticCount(deployedFlag *bool) (*ads.AdCount, error) {
	var resp *http.Response
	var err error
	if *deployedFlag {
		resp, err = http.Get("http://" + os.Getenv("elastic.host") + ":" + os.Getenv("elastic.port") + "/ads/ad/_count")

	} else {
		resp, err = http.Get("http://" + conf.Elastic.Host + ":" + conf.Elastic.Port + "/ads/ad/_count")

	}

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
func prepareBody(searchTerm string, filters map[string]string, fromSize map[string]string, priceRange string) string {
	//adList := &ads.AdList{}

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
		"lat",
		"lon"
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

	for k, v := range filters {
		terms = append(terms, fmt.Sprintf(`{"term": {"%v": %v}}`, k, v))
	}

	terms = append(terms, priceRange)

	body := fmt.Sprintf(search, fromSize["from"], fromSize["size"], queryParam, strings.Join(terms, ","))
	log.Printf("[%v] Query object sent to elastic:", remoteAddres)

	log.Println("\n", httputils.JSONPrettyPrint(body))
	return body
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

//GetAdListElastic this returns the ads.
//Pagination can be done using offset and limit
func GetAdListElastic(deployedFlag *bool, filter *ads.Filter) (*ads.AdList, error) {

	//this map will contain all the applicable filters received in the request
	//we must validate each type of filter to be able to set them properly for elasticSearch
	singleValueFilters := prepareSingleValueFilters(filter)
	fromSizeFilter := prepareFromSizeFilter(filter)
	priceRange := preparePriceRangeFilter(filter)

	prepareBody(filter.GetSearchParam().GetValue(), singleValueFilters, fromSizeFilter, priceRange)
	adList := &ads.AdList{}
	//get the results from elastic search
	//this needs to be changed for POST using the query parameters.
	// var prueba = `-d {"sort": [{ "field1": { "order": "desc" }},{ "field2": { "order": "desc" }}],"size": 100}`
	var err error
	var resp *http.Response
	if *deployedFlag {
		resp, err = http.Get("http://" + os.Getenv("elastic.host") + ":" + os.Getenv("elastic.port") + "/ads/ad/_search?from=" + strconv.Itoa(int(filter.From.GetValue())) + "&size=" + strconv.Itoa(int(filter.Size.GetValue())))
		log.Println("http://" + os.Getenv("elastic.host") + ":" + os.Getenv("elastic.port") + "/ads/ad/_search?from=" + strconv.Itoa(int(filter.From.GetValue())) + "m&size=" + strconv.Itoa(int(filter.Size.GetValue())))

	} else {

		resp, err = http.Get("http://" + conf.Elastic.Host + ":" + conf.Elastic.Port + "/ads/ad/_search?from=" + strconv.Itoa(int(filter.From.GetValue())) + "&size=" + strconv.Itoa(int(filter.Size.GetValue())))
		log.Println("http://" + conf.Elastic.Host + ":" + conf.Elastic.Port + "/ads/ad/_search?from=" + strconv.Itoa(int(filter.From.GetValue())) + "m&size=" + strconv.Itoa(int(filter.Size.GetValue())))

	}

	if err != nil {
		log.Println(err)
		return adList, err
	}

	if resp.StatusCode == 200 {
		// params := make(map[string]string)
		// params["sort"]=[{''}]
		// http.NewRequest("GET","http://" + conf.Elastic.Host + ":" + conf.Elastic.Port + "/ads/ad/_search",)

		defer resp.Body.Close()
		log.Println("Connected to Elastic...")
		//create a struct to hold the values from the response
		results := &elasticSearchAd{}

		body, err := ioutil.ReadAll(resp.Body)
		log.Println("Reading response from Elastic...")

		err = json.Unmarshal(body, &results)
		log.Println(fmt.Sprintf("took :%v, timed_out:%v ,hits: %v", results.Took, results.TimedOut, results.Hits.Total))
		if err != nil {
			panic(err)
		}

		adList.Ads = []*ads.Ad{}
		log.Println("Translating ads to protobuf...")
		//convert the ads to protobuf and add them to the adList that will be returned
		for _, ad := range results.Hits.Hits {
			adPB := &ads.Ad{}
			adPB = ToProto(ad.Source, adPB)
			adList.Ads = append(adList.Ads, adPB)

		}
		log.Println("done!")
		return adList, nil
	}
	//TODO: return custom errors like the ones coming from elastic, this will help troubleshoot in case of problems
	// {
	// 	"error": "Incorrect HTTP method for uri [/ads/ad/?sort] and method [GET], allowed: [POST]",
	// 	"status": 405
	// 	}
	return adList, errors.New("status" + strconv.Itoa(resp.StatusCode) + " MakakoLabs: There was a problem while procesing your request")

}

//SearchElastic serach in elastic search
func SearchElastic(deployedFlag *bool, filter *ads.Filter, remoteAddr string) (*ads.SearchResponse, error) {
	remoteAddres = remoteAddr
	//this map will contain all the applicable filters received in the request
	//we must validate each type of filter to be able to set them properly for elasticSearch
	myFilterMap := prepareSingleValueFilters(filter)
	fromSize := prepareFromSizeFilter(filter)
	priceRange := preparePriceRangeFilter(filter)

	requestBody := []byte(prepareBody(filter.GetSearchParam().GetValue(), myFilterMap, fromSize, priceRange))
	adList := &ads.AdList{}
	searchResponse := &ads.SearchResponse{}
	var err error
	var req *http.Request

	if *deployedFlag {

		req, err = http.NewRequest("POST", "http://"+os.Getenv("elastic_host")+":"+os.Getenv("elastic_port")+"/_search", bytes.NewBuffer(requestBody))
	} else {

		req, err = http.NewRequest("POST", "http://"+conf.Elastic.Host+":"+conf.Elastic.Port+"/_search", bytes.NewBuffer(requestBody))
	}

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
