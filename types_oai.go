package oaipmh

import (
	"encoding/xml"
	"time"
)

type ResponseError struct {
	XMLName xml.Name `xml:"OAI-PMH"`
	Error   Error    `xml:"error"`
}

type Error struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:",chardata"`
	Code    string   `xml:"code,attr"`
}

type MetadataFormat struct {
	XMLName           xml.Name `xml:"metadataFormat"`
	MetadataPrefix    string   `xml:"metadataPrefix"`
	Schema            string   `xml:"schema"`
	MetadataNamespace string   `xml:"metadataNamespace"`
}

type InterpretedRequest struct {
	XMLName xml.Name `xml:"request"`
	BaseURL string   `xml:",chardata"`
	Verb    string   `xml:"verb,attr"`
}

type ListMetadataFormatsOptions struct {
	Identifier string
}

type ListMetadataFormatsResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	ResponseDate       string             `xml:"responseDate"`
	MetadataFormats    []MetadataFormat   `xml:"ListMetadataFormats>metadataFormat"`
}

type Identify struct {
	XMLName           xml.Name `xml:"Identify"`
	RepositoryName    string   `xml:"repositoryName"`
	BaseURL           string   `xml:"baseURL"`
	ProtocolVersion   string   `xml:"protocolVersion"`
	EarliestDatestamp string   `xml:"earliestDatestamp"`
	DeletedRecord     string   `xml:"deletedRecord"`
	Granularity       string   `xml:"granularity"`
	AdminEmail        string   `xml:"adminEmail"`
	Compression       string   `xml:"compression"`
}

type IdentifyResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	ResponseDate       string             `xml:"responseDate"`
	Identify           Identify           `xml:"Identify"`
}

type RecordHeader struct {
	XMLName    xml.Name `xml:"header"`
	Identifier string   `xml:"identifier"`
	Datestamp  string   `xml:"datestamp"`
	SetSpec    []string `xml:"setSpec"`
	Status     string   `xml:"status"`
}

type Record struct {
	XMLName  xml.Name     `xml:"record"`
	Header   RecordHeader `xml:"header"`
	Metadata Metadata     `xml:"metadata"`
}

type Metadata struct {
	Raw []byte `xml:",innerxml"`
}

type GetRecordOptions struct {
	Identifier     string
	MetadataPrefix string
}

type GetRecordResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	Record             Record             `xml:"GetRecord>record"`
}

type ListOptions struct {
	MetadataPrefix  string
	From            time.Time
	Until           time.Time
	Set             string
	ResumptionToken string
}

type ListRecordsResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	Records            []Record           `xml:"ListRecords>record"`
	ResumptionToken    ResumptionToken    `xml:"ListRecords>resumptionToken"`
}

type ResumptionToken struct {
	XMLName        xml.Name `xml:"resumptionToken"`
	ExpirationDate string   `xml:"expirationDate,attr"`
	Value          string   `xml:",chardata"`
}

type ListIdentifiersResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	Headers            []RecordHeader     `xml:"ListIdentifiers>header"`
	ResumptionToken    ResumptionToken    `xml:"ListIdentifiers>resumptionToken"`
}

type ListSetsOptions struct {
	ResumptionToken string
}

type ListSetsResponse struct {
	XMLName            xml.Name           `xml:"OAI-PMH"`
	InterpretedRequest InterpretedRequest `xml:"request"`
	Sets               []Set              `xml:"ListSets>set"`
	ResumptionToken    ResumptionToken    `xml:"ListSets>resumptionToken"`
}

type Set struct {
	XMLName xml.Name `xml:"set"`
	SetSpec string   `xml:"setSpec"`
	SetName string   `xml:"setName"`
}
