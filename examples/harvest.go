package main

import (
	"github.com/nick-jones/oaipmh"
	"fmt"
	"os"
	"time"
)

func main() {
	baseURL := "http://eprints.ecs.soton.ac.uk/cgi/oai2"
	client, _ := oaipmh.NewClient(baseURL)
	options := &oaipmh.ListRecordsOptions{"oai_dc", time.Time{}, time.Time{}, "", ""}

	for {
		records := new(oaipmh.DublinCoreRecords)
		response, _, err := client.ListRecords(options, records)
		resumptionToken := response.ResumptionToken.Value
		options = &oaipmh.ListRecordsOptions{"", time.Time{}, time.Time{}, "", resumptionToken}

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, record := range records.Records {
			if len(record.Titles) > 0 {
				fmt.Fprintf(os.Stdout, "title: %s\n", record.Titles[0])
			}
		}

		if resumptionToken == "" {
			break
		}
	}
}