package main

import (
	"fmt"
	"net/http"
	"strconv"

	// TODO Improve and upstream own extension
	cloudant "github.com/IBM-Bluemix/go-cloudant"
	request "github.com/parnurzeal/gorequest"
)

// GetAirportRecords collects Airport records from Cloudant DB
func GetAirportRecords(request DistanceListRequest) ([]Airport, error) {
	airports := []Airport{}

	var (
		err      error
		client   *cloudant.Client
		db       *cloudant.DB
		design   *DesignDocumentExtended
		response *cloudant.SearchResp
	)

	query := fmt.Sprintf("lon:[%v TO %v] AND lat:[%v TO %v]",
		request.From.Lon, request.To.Lon,
		request.From.Lat, request.To.Lat,
	)

	if client, err = cloudant.NewClient(GlobalCloudantConfig.userName, ""); err != nil {
		return airports, err
	}
	db = client.DB(GlobalCloudantConfig.dbName)
	design = NewDesignDocumentExtended(GlobalCloudantConfig.designName)
	if response, err = design.SearchNotAuth(
		client, db, GlobalCloudantConfig.indexName, query, "", GlobalCloudantConfig.searchLimit,
	); err != nil {
		return airports, err
	}

	// TODO Move making Airport into separate function makeAirport(), []cloudant.searchRow must be visible
	// TODO Make []cloudant.searchRow visible ([]cloudant.SearchRow) and uptstream it
	for _, row := range response.Rows {
		var (
			fields map[string]interface{}
			ok     bool
			lat    float64
			lon    float64
		)

		if fields, ok = row.Fields.(map[string]interface{}); !ok {
			return airports, fmt.Errorf("unknown response format: %+v", row.Fields)
		}
		if lat, err = strconv.ParseFloat(fmt.Sprintf("%v", fields["lat"]), 64); err != nil {
			return airports, fmt.Errorf("invalid response format: %+v (%s)", row.Fields, err.Error())
		}
		if lon, err = strconv.ParseFloat(fmt.Sprintf("%v", fields["lon"]), 64); err != nil {
			return airports, fmt.Errorf("invalid response format: %+v (%s)", row.Fields, err.Error())
		}

		airport := Airport{
			ID:   row.ID,
			Name: fmt.Sprintf("%v", fields["name"]),
			Lat:  lat,
			Lon:  lon,
		}

		airports = append(airports, airport)
	}

	return airports, err
}

/*
 * Extension to cloudant.DesignDocument
 *
 * TODO Improve and upstream this extension to IBM-Bluemix/go-cloudant and/or cloudant-labs/go-cloudant
 */

// DesignDocumentExtended is an extension to cloudant.DesignDocument
type DesignDocumentExtended struct {
	cloudant.DesignDocument
}

// NewDesignDocumentExtended creates a new instance,
//    calls cloudant.NewDesignDocument
func NewDesignDocumentExtended(name string) *DesignDocumentExtended {
	ddoc := DesignDocumentExtended{
		DesignDocument: *cloudant.NewDesignDocument(name),
	}
	return &ddoc
}

// SearchNotAuth is an extended cloudant.DesignDocument.Search w/o authentication
// TODO Upstream extension to IBM-Bluemix/go-cloudant
// TODO Upstream similar to cloudant-labs/go-cloudant
// TODO Upstream increasing max. 200 limitation
func (ddoc *DesignDocumentExtended) SearchNotAuth(client *cloudant.Client, db *cloudant.DB,
	index, query, bookmark string, limit int) (*cloudant.SearchResp, error,
) {
	// TODO get dbPath from IBM-Bluemix/go-cloudant, after upstreaming (db.path is hidden now)
	dbPath := client.Client.URL() + "/" + db.Name()
	path := "/" + ddoc.ID + "/_search/" + index
	body := &cloudant.SearchResp{}
	req := request.New().
		Get(dbPath + path).
		Query("query=" + query).
		Query("limit=" + strconv.Itoa(limit))
	resp, respBody, errs := req.EndStruct(body)
	if errs != nil {
		return nil, errs[len(errs)-1]
	}
	// TODO upstream better error handling (in case of HTTP Status 4xx, the errs = nil)
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("cannot access Cloudant DB: %s; %s", resp.Status, string(respBody))
	}
	return body, nil
}
