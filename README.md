# airport-distance

Airport distance calculator

## Introduction

It's a simple example for calculating airport distances.

All examples are tested on Linux bash (Ubuntu 16.04, x86_64). Examples may work on Windows in Cygwin shells, for example:
* Cygwin shell
* Git bash
* MobaXterm local shell

## Restrictions

*We need to be able to review your code in a short time frame.*

TODO:

* Spit the source code to several directories (`cmd`, `util`, `common`, `api`, etc.) and files
* More strict and robust error handling (semantical check)

*Please use only a Cloudant client library and standard libraries.*

TODO:

* Using github.com/golang/glog for better logging
* Using github.com/spf13/{pflag,cobra,viper} for better parameter handling
* Using github.com/stretchr/testify for making unit tests
* More strict and robust error handling

* Airports (...) sorted by distance within a user provided radius of a user provided lat/long point.

TODO:

* Review ordering rule if distances are equal (different names at same lat/lon position)

## Cloudant client library

There is no official Go library at <https://cloud.ibm.com/docs/services/Cloudant/libraries?topic=cloudant-supported-client-libraries#su> and <https://cloud.ibm.com/docs/services/Cloudant?topic=cloudant-third-party-client-libraries>.

Finally, IBM-Cloud/go-cloudant was selected with a small extending part.

TODO:
* Checking if client libraries support proxy
* Solving and upstream increasing max. 200 Cloudant record limitation

Unofficial libraries:

### IBM-Cloud/go-cloudant

<https://github.com/IBM-Cloud/go-cloudant> is a simple wrapper to <https://github.com/timjacobi/go-couchdb>, builds `_design/view1/_search/geo` path postfix.

Does not support non-authentication (password is needed), see:

* `main_test.go:TestCloudant_IBM_Bluemix_go_cloudant()`
* <https://github.com/IBM-Cloud/go-cloudant/blob/master/cloudant.go#L213>

Last commit: Oct 13, 2016

TODO:
* Improve and upstream own extension, see `cloudant.go`, *Extension to cloudant.DesignDocument*

### timjacobi/go-couchdb

<https://github.com/timjacobi/go-couchdb> supports non-authenticated connection. Does not build `_design/view1/_search/geo` path postfix.

Last commit: Aug 17, 2016

### cloudant-labs/go-cloudant

<https://github.com/cloudant-labs/go-cloudant> does not build `_design/view1/_search/geo` path postfix. Skipped.

Last commit: Jun 20, 2018

TODO:
* Improve and upstream

### obieq/go-cloudant

<https://github.com/obieq/go-cloudant/> supports partly building `_design/view1/_search/geo` path postfix (design part, search part not). Skipped.

Last commit: Sep 29, 2014

## Building

Prerequisites:

* Go 1.11 or above (tested with `go1.11.5 and go1.12.7 linux/amd64`)
* Go modules environment (supported by Go 1.11 and above, see: `GO111MODULE` env variable)
* (Optional) Golint from github.com/golang/lint
* (Optional) Access to https://mikerhodes.cloudant.com for running tests

Source code can be download by below command:

```bash
git clone https://github.com/pgillich/airport-distance.git
```

Golint for optional static code checking can be installed by below command:

```bash
GO111MODULE=off go get -u golang.org/x/lint/golint
```

Running optional `go test` needs access to <https://mikerhodes.cloudant.com/airportdb/_design/view1>

### Building executable binary

Below command builds executable binary:

```bash
go build -mod=readonly
```

Optional: Below commands execute several checks:

```bash
go vet -mod=readonly
golint -min_confidence=.3
go test -mod readonly -v
```

### Building Docker image

Docker image build makes `go build` automatically. Example for building Docker image:

```bash
docker build -t pgillich/airport-distance .
```

## Running

All possible CLI parameters are documented in the built-in help, for example:

```bash
./airport-distance --help
```

### Running from CLI

Example for starting one operation:

```bash
./airport-distance -centerLon 3 -centerLat 2 -radius 500000
```

Coordinates are set in degree, radius in meter.

### Running as a REST service

Example for starting as a service:

```bash
./airport-distance -service
```

Exampe URL for getting an airport list:

```
http://localhost:8080/list/distance?centerLon=1&centerLat=1&radius=400000
```

### Running in a Docker container

```bash
docker run -it --rm -p 8080:8080 --name airport-distance pgillich/airport-distance
```

Exampe URL for getting an airport is same to above.

## Output format

The output/response is always JSON, if all CLI parameters can be parsed:
* `request`: the actual user parameters
* `errors`: a list of error messages, if any
* `airports`: a list of selected airports, sorted by distance (and name, if distance is same)

Distance is provided in meters.

Example for a response:

```bash
{
  "request": {
    "Radius": 500000,
    "Center": {
      "Lat": 2,
      "Lon": 3
    }
  },
  "errors": [],
  "airports": [
    {
      "id": "f6d336d85815d9ebc82fa5009fd73b03",
      "name": "Savut",
      "lat": 1,
      "lon": 1,
      "distance": 248568.719240913
    },
    {
      "id": "aa915feea5ecf2f87d0c7bca67cd6f4f",
      "name": "Xakan",
      "lat": 1,
      "lon": 1,
      "distance": 248568.719240913
    },
    {
      "id": "efb9627917c711e4e6852dfb7708172f",
      "name": "Xakan",
      "lat": 1,
      "lon": 1,
      "distance": 248568.719240913
    },
    {
      "id": "f6d336d85815d9ebc82fa5009fd7488e",
      "name": "Xugan",
      "lat": 1,
      "lon": 1,
      "distance": 248568.719240913
    },
    {
      "id": "aa915feea5ecf2f87d0c7bca6777b4f6",
      "name": "Les Ailerons",
      "lat": 0,
      "lon": 0,
      "distance": 400862.6244736863
    },
    {
      "id": "f6d336d85815d9ebc82fa5009fd12f34",
      "name": "Mainz Finthen",
      "lat": 0,
      "lon": 0,
      "distance": 400862.6244736863
    },
    {
      "id": "f6d336d85815d9ebc82fa5009fc56d04",
      "name": "Vilamendhoo",
      "lat": 0,
      "lon": 0,
      "distance": 400862.6244736863
    },
    {
      "id": "f6d336d85815d9ebc82fa5009f4b7181",
      "name": "Sao Tome Intl",
      "lat": 0.378175,
      "lon": 6.712153,
      "distance": 450353.71383098216
    },
(...)
```

Example for a failed response:

```bash
{
  "request": {
    "Radius": 500000,
    "Center": {
      "Lat": 1,
      "Lon": 1
    }
  },
  "errors": [
    "cannot access Cloudant DB: 404 Not Found; {\"error\":\"not_found\",\"reason\":\"BAD_INDEX not found.\"}\n"
  ],
  "airports": []
}
```
