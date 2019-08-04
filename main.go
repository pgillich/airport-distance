package main

// TODO Spit the source code to several directories (cmd, util, common, api, etc.) and files

import (
	// TODO Using github.com/spf13/{pflag,cobra,viper} for better parameter handling
	"flag"

	// TODO Using github.com/golang/glog for better logging

	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

// TODO move to common directory
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
	// OptSearchLimit is the limit of provided records by Cloudant
	OptSearchLimit = "searchLimit"

	// OptService starts REST service, if true
	OptService = "service"
	// OptListenOn sets host:port where the service listens on
	OptListenOn = "listenOn"

	// RadDeg is 1 rad in degree: 57.295779513
	RadDeg = 180.0 / math.Pi
	// Earth radius in meter
	EarthRadius = 6371e3
)

// CloudantConfig contains Coudant client and DB config
type CloudantConfig struct {
	userName    string
	dbName      string
	designName  string
	indexName   string
	searchLimit int
}

// ServiceConfig describes the REST service config
type ServiceConfig struct {
	startService bool
	listenOn     string
}

// TODO Move to common directory
var (
	GlobalCliRequest     DistanceListRequest
	GlobalCloudantConfig CloudantConfig
	GlobalServiceConfig  ServiceConfig
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

// Airport record type (REST reponse)
type Airport struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Distance float64 `json:"distance"`
}

// DistanceListResponse contains response for the CLI or REST query
// TODO replace Errors to []error and implementing custom marshaller
type DistanceListResponse struct {
	Request  DistanceListRequest `json:"request"`
	Errors   []string            `json:"errors"`
	Airports []Airport           `json:"airports"`
}

func main() {
	if GlobalServiceConfig.startService {
		startService()
	} else {
		executeCli()
	}
}

// startServiceHandler start the REST service
func startService() {
	log.Printf("Listening on %s\n", GlobalServiceConfig.listenOn)

	err := http.ListenAndServe(GlobalServiceConfig.listenOn, nil)
	if err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}

// handleDistanceList is the handler function of /list/distance
// TODO Remove JSON indentation
func handleDistanceList(w http.ResponseWriter, r *http.Request) {
	responseBytes := []byte{}
	w.Header().Set("Content-Type", "application/json")

	parserErrors := make([]string, 0, 6) // capacity: max. number of parser error
	values := r.URL.Query()

	request := DistanceListRequest{
		From: GeoCoord{
			Lat: parseFloatQueryParam(values, OptFromLat, GlobalCliRequest.From.Lat, &parserErrors),
			Lon: parseFloatQueryParam(values, OptFromLon, GlobalCliRequest.From.Lon, &parserErrors),
		},
		To: GeoCoord{
			Lat: parseFloatQueryParam(values, OptToLat, GlobalCliRequest.To.Lat, &parserErrors),
			Lon: parseFloatQueryParam(values, OptToLon, GlobalCliRequest.To.Lon, &parserErrors),
		},
		Center: GeoCoord{
			Lat: parseFloatQueryParam(values, OptCenterLat, GlobalCliRequest.Center.Lat, &parserErrors),
			Lon: parseFloatQueryParam(values, OptCenterLon, GlobalCliRequest.Center.Lon, &parserErrors),
		},
	}

	if len(parserErrors) > 0 {
		response := DistanceListResponse{
			Request:  request,
			Errors:   parserErrors,
			Airports: []Airport{},
		}
		// TODO handling marshalling error
		responseBytes, _ = json.MarshalIndent(response, "", "  ")
		w.WriteHeader(http.StatusPreconditionFailed)
	} else {
		response, err := doDistanceList(GlobalCliRequest)
		// TODO handling marshalling error
		responseBytes, _ = json.MarshalIndent(response, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
		}
	}

	w.Write(responseBytes)
}

// parseFloatQueryParam parses a float and adds error message to errors, if needed
func parseFloatQueryParam(values url.Values, name string, defaultValue float64, errors *[]string) float64 {
	value := defaultValue
	var err error

	if valueStr := values.Get(name); valueStr != "" {
		if value, err = strconv.ParseFloat(valueStr, 64); err != nil {
			(*errors) = append(*errors, err.Error())
		}
	}

	return value
}

// executeCli makes 1 CLI operation
// TODO Remove JSON indentation
func executeCli() {
	var err error

	response, err := doDistanceList(GlobalCliRequest)
	// TODO handling marshalling error
	responseBytes, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(responseBytes))
	if err != nil {
		os.Exit(1)
	}
}

// doDistanceList is the main function to collect, order and marshall response
// The response is structure, which contains error messages, too; ready to make a JSON text
func doDistanceList(request DistanceListRequest) (DistanceListResponse, error) {
	var err error

	reponse := DistanceListResponse{
		Request:  request,
		Errors:   []string{},
		Airports: []Airport{},
	}

	if err = checkCorrectDistanceListRequest(request); err != nil {
		reponse.Errors = append(reponse.Errors, err.Error())
	} else if reponse.Airports, err = GetAirportRecords(request); err != nil {
		reponse.Errors = append(reponse.Errors, err.Error())
	} else {
		calculateDistances(reponse.Airports, GlobalCliRequest.Center)
		orderByDistance(reponse.Airports)
	}

	return reponse, err
}

// checkCorrectDistanceListRequest checks and corrects a DistanceListRequest
// TODO More strict and robust error handling (semantical check)
func checkCorrectDistanceListRequest(request DistanceListRequest) error {
	return nil
}

// calculateDistances calculates distance between the airports and center coordinates
func calculateDistances(airports []Airport, center GeoCoord) {
	for a, airport := range airports {
		airports[a].Distance = calculateDistance(
			airport.Lat/RadDeg, airport.Lon/RadDeg,
			center.Lat/RadDeg, center.Lon/RadDeg,
		)
	}
}

// calculateDistance calculates distance between 2 spherical coords (input in radians)
// See: https://www.movable-type.co.uk/scripts/latlong.html
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	R := EarthRadius

	fi1 := lat1
	fi2 := lat2
	dFi := lat2 - lat1
	dLambda := lon2 - lon1

	a := math.Sin(dFi/2.0)*math.Sin(dFi/2.0) +
		math.Cos(fi1)*math.Cos(fi2)*math.Sin(dLambda/2.0)*math.Sin(dLambda/2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// orderByDistance sorts airports in ascending distance order
// TODO Review ordering rule if distances are equal (different names at same lat/lon position)
func orderByDistance(airports []Airport) {
	sort.Slice(airports, func(i, j int) bool {
		if airports[i].Distance == airports[j].Distance {
			return strings.Compare(airports[i].Name, airports[j].Name) < 0
		}
		return airports[i].Distance < airports[j].Distance
	})
}

// parseParams is called by init()
func parseParams() {
	flag.Float64Var(&GlobalCliRequest.From.Lat, OptFromLat, 0.0, "Latitude of From corner")
	flag.Float64Var(&GlobalCliRequest.From.Lon, OptFromLon, 0, "Longitude of From corner")
	flag.Float64Var(&GlobalCliRequest.To.Lat, OptToLat, 30.0, "Latitude of To corner")
	flag.Float64Var(&GlobalCliRequest.To.Lon, OptToLon, 5.0, "Longitude of To corner")
	flag.Float64Var(&GlobalCliRequest.Center.Lat, OptCenterLat, 2.0, "Latitude of Center corner")
	flag.Float64Var(&GlobalCliRequest.Center.Lon, OptCenterLon, 3.0, "Longitude of Center corner")

	flag.StringVar(&GlobalCloudantConfig.userName, OptUserName, "mikerhodes", "Cloudant user name")
	flag.StringVar(&GlobalCloudantConfig.dbName, OptDbName, "airportdb", "Cloudant DB name")
	flag.StringVar(&GlobalCloudantConfig.designName, OptDesignName, "view1", "Cloudant design name")
	flag.StringVar(&GlobalCloudantConfig.indexName, OptIndexName, "geo", "Cloudant index name")
	flag.IntVar(&GlobalCloudantConfig.searchLimit, OptSearchLimit, 100, "Limit of provided records by Cloudant, must not exceed 200")

	flag.BoolVar(&GlobalServiceConfig.startService, OptService, false, "If true: starts REST service, instead of CLI operation")
	flag.StringVar(&GlobalServiceConfig.listenOn, OptListenOn, ":8080", "host:port where the REST service listens on")

	flag.Parse()
}

func init() {
	parseParams()
	http.HandleFunc("/list/distance", handleDistanceList)
}
