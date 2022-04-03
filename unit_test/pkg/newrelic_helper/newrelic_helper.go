package newrelic_helper

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

func InitNewRelic(licenseKey, newRelicAppName string, newRelicLicenseKey string) *newrelic.Application {

	var err error
	newRelicApp, err := newrelic.NewApplication(newrelic.ConfigAppName(newRelicAppName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true))
	if err != nil {
		panic("Failed to setup NewRelic: " + err.Error())
	}

	txn := newRelicApp.StartTransaction("Init NewRelic")

	defer txn.End()

	s := newrelic.DatastoreSegment{
		Product:            newrelic.DatastoreCassandra,
		Collection:         "book",
		Operation:          "GET",
		ParameterizedQuery: "SELECT * FROM book",
		QueryParameters: map[string]interface{}{
			"id": 5,
		},
		DatabaseName: "book_store",
	}
	s.StartTime = txn.StartSegmentNow()

	s.End()

	return newRelicApp
}
