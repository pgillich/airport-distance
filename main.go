package main

import (
	"flag"
	"log"
)

const (
	// OptFromLat is the latitude of From corner
	OptFromLat = "fromLat"
	// OptFromLon is the longitude of From corner
	OptFromLon = "fromLon"
	// OptToLat is the latitude of To corner
	OptToLat = "toLat"
	// OptToLon is the longitude of From corner
	OptToLon = "toLon"
	// OptCenterLat is the latitude of Center (a point of the distance)
	OptCenterLat = "centerLat"
	// OptCenterLon is the longitude of Center (a point of the distance)
	OptCenterLon = "centerLon"
	// OptUserName is the Cloudant user name
	OptUserName = "userName"
	// OptDbName is the Cloudant DB name
	OptDbName = "dbName"
	// OptDesignName is the Cloudant design name
	OptDesignName = "designName"
	// OptIndexName is the Cloudant index name
	OptIndexName = "indexName"

	// OptService starts REST service, if true
	OptService = "service"
)

var (
	cliRequest DistanceListRequest

	userName   string
	dbName     string
	designName string
	indexName  string

	startService bool
)

// GeoCoord is a 2D coordinate type (latitude, longitude)
type GeoCoord struct {
	Lat float64
	Lon float64
}

// DistanceListRequest contains the CLI or REST parameters for distance list calculation
type DistanceListRequest struct {
	From   GeoCoord
	To     GeoCoord
	Center GeoCoord
}

// Airport record type
type Airport struct {
	ID       string
	Name     string
	Lat      float64
	Lon      float64
	Distance float64
}

func main() {
	var err error

	if startService {
		startServiceHandler()
	} else {
		var airports []Airport

		if err = checkDistanceListRequest(cliRequest); err != nil {
			log.Fatalf("Invalid parameter: %+v\n", err.Error())
		}
		if airports, err = getAirportRecords(cliRequest); err != nil {
			log.Fatalf("Cannot get data: %+v\n", err.Error())
		}
		calculateDistances(airports, cliRequest.Center)
	}
}

func checkDistanceListRequest(request DistanceListRequest) error {
	return nil
}

func calculateDistances(airports []Airport, center GeoCoord) {

}

func startServiceHandler() {

}

func parseParams() {
	flag.Float64Var(&cliRequest.From.Lat, OptFromLat, 0.0, "Latitude of From corner")
	flag.Float64Var(&cliRequest.From.Lon, OptFromLon, 0, "Longitude of From corner")
	flag.Float64Var(&cliRequest.To.Lat, OptToLat, 30.0, "Latitude of To corner")
	flag.Float64Var(&cliRequest.To.Lon, OptToLon, 5.0, "Longitude of To corner")
	flag.Float64Var(&cliRequest.Center.Lat, OptCenterLat, 2.0, "Latitude of Center corner")
	flag.Float64Var(&cliRequest.Center.Lon, OptCenterLon, 3.0, "Longitude of Center corner")
	flag.StringVar(&userName, OptUserName, "mikerhodes", "Cloudant user name")
	flag.StringVar(&dbName, OptDbName, "airportdb", "Cloudant DB name")
	flag.StringVar(&designName, OptDesignName, "view1", "Cloudant design name")
	flag.StringVar(&indexName, OptIndexName, "geo", "Cloudant index name")
	flag.BoolVar(&startService, OptService, false, "If true: starts REST service, instead of CLI operation, all Latitude/Longitude params are skipped")
	flag.Parse()
}
