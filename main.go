package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/jebo87/makako-api/repository"
	"gitlab.com/jebo87/makako-gateway/httputils"
	"gitlab.com/jebo87/makako-grpc/ads"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type adsServer struct{}

var deployedFlag *bool

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	deployedFlag = flag.Bool("deployed", false, "Defines if absolute paths need to be used for the config files")

	flag.Parse()
	repository.InitializeDBConfig(deployedFlag)

	// create a listener on TCP port 7777
	var listener net.Listener
	var err error
	if *deployedFlag {
		listener, err = net.Listen("tcp", os.Getenv("PROD_ADDRESS")+":7777")

	} else {
		listener, err = net.Listen("tcp", os.Getenv("DEV_ADDRESS")+":7777")
		log.Println(os.Getenv("DEV_ADDRESS") + ":7777")

	}

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
	adFromDB, err := repository.GetAdPB(text.Text)
	log.Println(adFromDB.GetTitle())
	log.Println("AdDetail: Sending response")
	return adFromDB, err

}

func (adsServer) Count(ctx context.Context, void *ads.Void) (count *ads.AdCount, err error) {
	log.Println("AdDetail: gRPC connection for ad count")
	count, err = repository.GetElasticCount(deployedFlag)
	log.Println("AdDetail: Sending response")
	return count, err

}
func (adsServer) AddListing(ctx context.Context, listing *ads.Ad) (count *ads.ListingID, err error) {
	log.Println("AddListing: gRPC call started")
	log.Println("AdDetail: Sending response")
	return &ads.ListingID{ListingID: 182919}, nil

}
func (adsServer) List(ctx context.Context, filter *ads.Filter) (*ads.SearchResponse, error) {

	httputils.LogDivider()
	peerInfo, _ := peer.FromContext(ctx)
	log.Printf("[%v] Remote gRPC Client", peerInfo.Addr)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("%v Procesing request from remote address", md["remote-addr"])
	//from database:
	// ads, err := repository.GetAdListPB(0, 0)

	//from elastic search
	ads, err := repository.SearchElastic(deployedFlag, filter, md["remote-addr"][0])

	if err != nil {
		log.Printf("%v Error %v", md["remote-addr"], err)

	}

	// ads, err := repository.GetAdListElastic(filter)
	// log.Println("printing from List in main:")
	// log.Println(ads.Ads)
	log.Printf("%v Finished processing request ", md["remote-addr"])

	return ads, err

}
