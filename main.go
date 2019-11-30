package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"bitbucket.org/jebo87/makako-api/store"
	"bitbucket.org/jebo87/makako-grpc/ads"

	"google.golang.org/grpc"
)

type adsServer struct{}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	store.InitializeDB()

	// create a listener on TCP port 7777
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	var adsSer adsServer

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	log.Println("Creating gRPC server...")
	// attach the ads service to the server
	ads.RegisterAdsServer(grpcServer, &adsSer)
	log.Println("Attaching ads service..")

	// start the server
	log.Println("Serving and waiting for connections in port 7777...")
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	<-c
}

func (adsServer) AdDetail(ctx context.Context, text *ads.Text) (ad *ads.Ad, err error) {
	log.Println("AdDetail: gRPC connection for adDetail ID: ", text.Text)
	adFromDB, err := store.GetAdPB(text.Text)
	log.Println(adFromDB.GetTitle())
	log.Println("AdDetail: Sending response")
	return adFromDB, err

}

func (adsServer) List(ctx context.Context, void *ads.Void) (*ads.AdList, error) {
	log.Println("List: loading ads..")
	//from database:
	// ads, err := store.GetAdListPB(0, 0)

	//from elastic search
	ads, err := store.GetAdListElastic(0, 0)
	// log.Println("printing from List in main:")
	// log.Println(ads.Ads)
	log.Println("List: Ads loaded ")
	return ads, err

}

// //ServeHTTP for the App
// func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {

// 	var head string
// 	head, req.URL.Path = ShiftPath(req.URL.Path)
// 	if head == "ads" {

// 		h.AdsHandler.ServeHTTP(res, req)
// 		return
// 	}
// 	http.Error(res, "Not Found", http.StatusNotFound)
// }

// //ServeHTTP for the Ads
// func (h *AdsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	enableCors(&res)
// 	var head string
// 	head, _ = ShiftPath(req.URL.Path)
// 	//validate if there is an actual id
// 	if id, _ := strconv.Atoi(head); id != 0 {
// 		//if there is and ID then the AdHandler
// 		//should take care of bringing that specific ad
// 		h.AdHandler.ServeHTTP(res, req)
// 		return
// 	}

// 	//check if there is an offset and a limit in the query parameters.
// 	offset, errOffset := strconv.Atoi(req.URL.Query().Get("offset"))
// 	limit, errLimit := strconv.Atoi(req.URL.Query().Get("limit"))

// 	//default to zero if offset or limit are not set
// 	if errOffset != nil {
// 		offset = 0
// 	}
// 	if errLimit != nil {
// 		limit = 0
// 	}

// 	switch req.Method {
// 	case "GET":
// 		fmt.Println("loading ads, request from " + req.RemoteAddr)
// 		res.Write(store.GetAdList(offset, limit))
// 	default:
// 		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)

// 	}
// 	return

// }

// //ServeHTTP for one Ad
// func (h *AdHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	var head string
// 	head, req.URL.Path = ShiftPath(req.URL.Path)
// 	switch req.Method {
// 	case "GET":
// 		fmt.Println("loading ad " + head)
// 		res.Write(store.GetAd(head))
// 	default:
// 		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)

// 	}
// 	return
// }

// //ShiftPath returns the head of the URL without initial slash '/' and the rest of the URL
// func ShiftPath(p string) (head, tail string) {
// 	p = path.Clean("/" + p)
// 	i := strings.Index(p[1:], "/") + 1
// 	if i <= 0 {
// 		return p[1:], "/"
// 	}
// 	return p[1:i], p[i:]
// }

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
// }

// func testAddAd() {
// 	var ad store.Ad
// 	ad.Title = "prueba 1"
// 	ad.Description = "Description 1"
// 	ad.City = "Montreal"
// 	ad.Country = "Canada"
// 	ad.Price = 660
// 	ad.PublishedDate = time.Now()
// 	ad.PropertyType = "apartment"
// 	ad.Rooms = 4
// 	ad.UserAdID = 1

// 	store.InsertAd(ad)

// }
