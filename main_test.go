package main

import (
	"fmt"
	"os"
	"testing"

	cloudant "github.com/IBM-Bluemix/go-cloudant"
)

func TestMain(m *testing.M) {
	parseParams()
	os.Exit(m.Run())
}

/*
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

func TestCloudant_IBM_Bluemix_go_cloudant_Extended(t *testing.T) {
	var err error

	client, err := cloudant.NewClient("mikerhodes", "")
	fmt.Printf("Client: %+v, %+v\n", client, err)

	db := client.DB("airportdb")
	fmt.Printf("DB: %+v, %+v\n", db, err)

	design := NewDesignDocumentExtended("view1")
	fmt.Printf("Design: %+v\n", design)

	result, err := design.SearchNotAuth(client, db, "geo", "lon:[0 TO 30] AND lat:[0 TO 5]")
	fmt.Printf("Result: %+v, %+v\n", result, err)

	return
}

func TestCloudant_getRecords_1(t *testing.T) {
	var err error

	request := DistanceListRequest{
		From:   GeoCoord{Lat: 0.0, Lon: 0.0},
		To:     GeoCoord{Lat: 5.0, Lon: 30.0},
		Center: GeoCoord{Lat: 2.0, Lon: 3.0},
	}

	airports, err := getAirportRecords(request)
	fmt.Printf("Airports: %+v, %+v\n", airports, err)
}
