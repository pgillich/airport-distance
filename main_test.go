package main

import (
	"fmt"
	"math"
	"testing"

	// TODO github.com/stretchr/testify for making unit tests (value checks are skipped now)

	cloudant "github.com/IBM-Bluemix/go-cloudant"
)

// TODO Faking Cloudant responses for unit tests, see: https://blog.pragmatists.com/test-doubles-fakes-mocks-and-stubs-1a7491dfa3da

func TestCloudant_IBM_Bluemix_go_cloudant_Extended(t *testing.T) {
	var err error

	client, err := cloudant.NewClient("mikerhodes", "")
	//fmt.Printf("Client: %+v, %+v\n", client, err)
	if err != nil {
		t.Fatalf("Cannot create client: %+v", err)
	}

	db := client.DB("airportdb")
	//fmt.Printf("DB: %+v, %+v\n", db, err)

	design := NewDesignDocumentExtended("view1")
	//fmt.Printf("Design: %+v\n", design)

	result, err := design.SearchNotAuth(client, db, "geo", "lon:[0 TO 30] AND lat:[0 TO 5]", "", GlobalCloudantConfig.searchLimit)
	//fmt.Printf("Result: %+v, %+v\n", result, err)
	if err != nil {
		t.Fatalf("Cannot get response: %+v %+v", result, err)
	}

	request := DistanceListRequest{
		Radius: 222000.0,
		Center: GeoCoord{Lat: 2.0, Lon: 3.0},
	}

	airports, err := GetAirportRecords(request)
	if err != nil {
		fmt.Printf("Airports: %+v\n", airports)
		t.Fatalf("Cannot get airports: %+v", err)
	}

	calculateDistances(airports, request.Center)
	filteredAirports := filterByRadius(airports, request.Radius)
	orderByDistance(filteredAirports)

	return
}

func TestDistance_doDistanceList(t *testing.T) {
	request := DistanceListRequest{
		Radius: 333000.0,
		Center: GeoCoord{Lat: 2.0, Lon: 3.0},
	}

	response, err := doDistanceList(request)
	if err != nil {
		fmt.Printf("Response: %+v\n", response)
		t.Fatalf("Error: %+v", err)
	} else {
		//fmt.Printf("Response: %+v\n", response)
	}
}

func TestDistance_calculateDistance_table(t *testing.T) {
	// See: https://www.movable-type.co.uk/scripts/latlong.html?from=64.1265,-21.8174&to=40.7128,-74.0060
	distanceTolerance := 1.0

	table := []struct {
		a GeoCoord
		b GeoCoord
		d float64
	}{
		{
			GeoCoord{64.1265, -21.8174},
			GeoCoord{40.7128, -74.0060},
			4.208198758424172e+06,
		},
		{
			GeoCoord{64.1265, 21.8174},
			GeoCoord{40.7128, 74.0060},
			4.208198758424172e+06,
		},
		{
			GeoCoord{-64.1265, 21.8174},
			GeoCoord{40.7128, 74.0060},
			1.2519184867839511e+07,
		},
		{
			GeoCoord{64.1265, -21.8174},
			GeoCoord{40.7128, 74.0060},
			6.271918839709706e+06,
		},
		{
			GeoCoord{-64.1265, -21.8174},
			GeoCoord{40.7128, 74.0060},
			1.4271722096567325e+07,
		},
		{
			GeoCoord{64.1265, 21.8174},
			GeoCoord{-40.7128, 74.0060},
			1.2519184867839511e+07,
		},
		{
			GeoCoord{64.1265, 21.8174},
			GeoCoord{40.7128, -74.0060},
			6.271918839709706e+06,
		},
		{
			GeoCoord{64.1265, 21.8174},
			GeoCoord{-40.7128, -74.0060},
			1.4271722096567325e+07,
		},
		{
			GeoCoord{0, 0},
			GeoCoord{0, 0},
			0.0,
		},
	}

	for _, testCase := range table {
		d := calculateDistance(testCase.a.Lat/RadDeg, testCase.a.Lon/RadDeg, testCase.b.Lat/RadDeg, testCase.b.Lon/RadDeg)
		if math.Abs(testCase.d-d) > distanceTolerance {
			t.Errorf("Invalid distance between %+v; %+v: %v != %v", testCase.a, testCase.b, d, testCase.d)
		}
	}

	return
}

func deltaAbs(a, b float64) float64 {
	return math.Abs(a - b)
}

/*
// TestCloudant_IBM_Bluemix_go_cloudant is failed code, because IBM-Bluemix/go-cloudant doesn't support non-auth
func TestCloudant_IBM_Bluemix_go_cloudant(t *testing.T) {
	var err error

	client, err := cloudant.NewClient("mikerhodes", "")
	fmt.Printf("Client: %+v, %+v\n", client, err)

	// Returns error:
	// PUT https://mikerhodes.cloudant.com/airportdb: (401) unauthorized: Name or password is incorrect.
	//
	//db, err := client.EnsureDB("airportdb")
	//fmt.Printf("DB: %+v, %+v\n", db, err)

	db := client.DB("airportdb")
	fmt.Printf("DB: %+v, %+v\n", db, err)

	design := cloudant.NewDesignDocument("view1")
	fmt.Printf("Design: %+v\n", design)

	result, err := design.Search(db, "geo", "lon:[0 TO 30] AND lat:[0 TO 5]", "", 0)
	fmt.Printf("Result: %+v, %+v\n", result, err)

	return
}
*/

/*
// TestCloudant_timjacobi_go_couchdb is an example for non-auth connection
func TestCloudant_timjacobi_go_couchdb(t *testing.T) {
	var err error

	username := "mikerhodes"
	url := fmt.Sprintf("https://%s.cloudant.com", username)
	client, err := couchdb.NewClient(url, nil)
	fmt.Printf("Client: %+v, %+v\n", client, err)

	err = client.Ping()
	fmt.Printf("Ping: %+v\n", err)

	return
}
*/
