package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"gitlab.com/jebo87/makako-api/store"
	"gitlab.com/jebo87/makako-grpc/ads"

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

func (adsServer) Count(ctx context.Context, void *ads.Void) (count *ads.AdCount, err error) {
	log.Println("AdDetail: gRPC connection for ad count")
	count, err = store.GetElasticCount()
	log.Println("AdDetail: Sending response")
	return count, err

}
func (adsServer) List(ctx context.Context, filter *ads.Filter) (*ads.AdList, error) {
	log.Println("List: loading ads..")

	//from database:
	// ads, err := store.GetAdListPB(0, 0)

	//from elastic search
	ads, err := store.GetAdListElastic(int(filter.From), int(filter.Size))
	// log.Println("printing from List in main:")
	// log.Println(ads.Ads)
	log.Println("List: Ads loaded ")
	return ads, err

}
