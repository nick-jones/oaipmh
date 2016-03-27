package oaipmh

import "encoding/xml"

type DublinCoreRecords struct {
	Records []DublinCoreRecord
}

type DublinCoreRecord struct {
	XMLName      xml.Name `xml:"http://www.openarchives.org/OAI/2.0/oai_dc/ dc"`
	Titles       []string `xml:"http://purl.org/dc/elements/1.1/ title"`
	Creators     []string `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Subjects     []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
	Descriptions []string `xml:"http://purl.org/dc/elements/1.1/ description"`
	Publishers   []string `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Contributors []string `xml:"http://purl.org/dc/elements/1.1/ contributor"`
	Dates        []string `xml:"http://purl.org/dc/elements/1.1/ date"`
	Types        []string `xml:"http://purl.org/dc/elements/1.1/ type"`
	Formats      []string `xml:"http://purl.org/dc/elements/1.1/ format"`
	Identifiers  []string `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Sources      []string `xml:"http://purl.org/dc/elements/1.1/ source"`
	Languages    []string `xml:"http://purl.org/dc/elements/1.1/ language"`
	Relations    []string `xml:"http://purl.org/dc/elements/1.1/ relation"`
	Coverages    []string `xml:"http://purl.org/dc/elements/1.1/ coverage"`
	Rights       []string `xml:"http://purl.org/dc/elements/1.1/ rights"`
}
