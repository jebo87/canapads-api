package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/jebo87/makako-api/store"
)

//App struct
type App struct {
	AdsHandler *AdsHandler
}

//AdsHandler is the handler for all ads requests
type AdsHandler struct {
}

func main() {

	//check if port was set, if not default to 8081
	var (
		port = flag.String("port", "", "payload data")
	)
	flag.Parse()

	if *port == "" {
		*port = "8081"
	}
	fmt.Println("Loading makako API server...")
	fmt.Println("Listening on port " + *port + " ...")

	//initializes the connection to the database
	store.InitializeDB()

	app := &App{
		AdsHandler: new(AdsHandler),
	}

	http.ListenAndServe(":"+*port, app)

	// testAddAd()
	//store.DeleteAd(9)

}

//ServeHTTP for the App
func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	if head == "ads" {

		h.AdsHandler.ServeHTTP(res, req)
		return
	}
	http.Error(res, "Not Found", http.StatusNotFound)
}

//ServeHTTP for the Ads
func (h *AdsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	offset, errOffset := strconv.Atoi(req.URL.Query().Get("offset"))
	limit, errLimit := strconv.Atoi(req.URL.Query().Get("limit"))

	if errOffset != nil {
		offset = 0
	}
	if errLimit != nil {
		limit = 0
	}

	switch req.Method {
	case "GET":
		res.Write(store.GetAdTitles(offset, limit))
	default:
		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)

	}
	return

}

//ShiftPath returns the head of the URL without initial slash '/' and the rest of the URL
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
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
