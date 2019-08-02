# airport-distance

Airport distance calculator

## Introduction

It's a simple example for calculating airport distances.

## Restrictions

*We need to be able to review your code in a short time frame.*

TODO:

* Spit the source code to several directories (`cmd`, `util`, `common`, `api`, etc.) and files

*Please use only a Cloudant client library and standard libraries.*

TODO:

* Using github.com/golang/glog for better logging
* Using github.com/spf13/{pflag,cobra,viper} for better parameter handling
* Using github.com/stretchr/testify for unit tests
* Robust error handling

## Cloudant client library

There is no official Go library at <https://cloud.ibm.com/docs/services/Cloudant/libraries?topic=cloudant-supported-client-libraries#su> and <https://cloud.ibm.com/docs/services/Cloudant?topic=cloudant-third-party-client-libraries>. 

Finally, IBM-Cloud/go-cloudant was selected and a small extended part.

TODO: extending IBM-Cloud/go-cloudant with custom authentication (no auth)

Unofficial libraries:

### IBM-Cloud/go-cloudant

<https://github.com/IBM-Cloud/go-cloudant> is a simple wrapper to <https://github.com/timjacobi/go-couchdb>.

Does not support non-authentication (password is needed), see:

* `main_test.go:TestCloudant_IBM_Bluemix_go_cloudant()`
* <https://github.com/IBM-Cloud/go-cloudant/blob/master/cloudant.go#L213>

Last commit: Oct 13, 2016

### timjacobi/go-couchdb

<https://github.com/timjacobi/go-couchdb> supports non-authenticated connection. Does not build `_design/view1/_search/geo` path postfix.

Last commit: Aug 17, 2016

### cloudant-labs/go-cloudant

<https://github.com/cloudant-labs/go-cloudant> does not build `_design/view1/_search/geo` path postfix. Skipped.

Last commit: Jun 20, 2018

### obieq/go-cloudant

<https://github.com/obieq/go-cloudant/> supports partly building `_design/view1/_search/geo` path postfix (desig part, search part not). Skipped.

Last commit: Sep 29, 2014

## Building

Prerequisites:

* Go 1.11 or above (tested with `go1.11.5 linux/amd64`)
* Go modules environment (supported by Go 1.11 and above)
* (optional) Golint from github.com/golang/lint

