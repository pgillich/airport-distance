package main

import (
	"fmt"
	"strconv"

	cloudant "github.com/IBM-Bluemix/go-cloudant"
	request "github.com/parnurzeal/gorequest"
)

func getAirportRecords(request DistanceListRequest) ([]Airport, error) {
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

	if client, err = cloudant.NewClient(userName, ""); err != nil {
		return airports, err
	}
	db = client.DB(dbName)
	design = NewDesignDocumentExtended(designName)
	if response, err = design.SearchNotAuth(client, db, indexName, query); err != nil {
		return airports, err
	}

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
 */

// DesignDocumentExtended is an extension to cloudant.DesignDocument
// TODO: Upstream extension to IBM-Bluemix/go-cloudant
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
// TODO: Upstream extension to IBM-Bluemix/go-cloudant
func (ddoc *DesignDocumentExtended) SearchNotAuth(client *cloudant.Client, db *cloudant.DB,
	index, query string) (*cloudant.SearchResp, error,
) {
	dbPath := client.Client.URL() + "/" + db.Name()
	path := "/" + ddoc.ID + "/_search/" + index
	body := &cloudant.SearchResp{}
	req := request.New().
		Get(dbPath + path).
		Query("query=" + query)
	if _, _, errs := req.EndStruct(body); errs != nil {
		return nil, errs[len(errs)-1]
	}
	return body, nil
}
