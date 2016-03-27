package oaipmh

import (
	"encoding/xml"
	"fmt"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

// Adapted from http://keighl.com/post/mocking-http-responses-in-golang/
func mockClient(code int, body string) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	http := &http.Client{Transport: tr}

	client := &Client{
		baseURL: server.URL,
		http:    http,
	}

	return server, client
}

func (s *S) TestListMetadataFormats(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T18:17:43Z</responseDate>
  <request verb="ListMetadataFormats">http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <ListMetadataFormats>
    <metadataFormat>
      <metadataPrefix>oai_bibl</metadataPrefix>
      <schema>http://www.openarchives.org/OAI/2.0/oai_dc.xsd</schema>
      <metadataNamespace>http://www.openarchives.org/OAI/2.0/oai_dc/</metadataNamespace>
    </metadataFormat>
    <metadataFormat>
      <metadataPrefix>oai_dc</metadataPrefix>
      <schema>http://www.openarchives.org/OAI/2.0/oai_dc.xsd</schema>
      <metadataNamespace>http://www.openarchives.org/OAI/2.0/oai_dc/</metadataNamespace>
    </metadataFormat>
  </ListMetadataFormats>
</OAI-PMH>
	`)

	defer server.Close()
	response, err := client.ListMetadataFormats("")

	expected := &ListMetadataFormatsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "ListMetadataFormats",
		},
		Error:        ResponseError{XMLName: xml.Name{Space: "", Local: ""}, Message: "", Code: ""},
		ResponseDate: "2016-03-26T18:17:43Z",
		MetadataFormats: []MetadataFormat{
			MetadataFormat{
				XMLName:           xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "metadataFormat"},
				MetadataPrefix:    "oai_bibl",
				Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
				MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/"},
			MetadataFormat{
				XMLName:           xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "metadataFormat"},
				MetadataPrefix:    "oai_dc",
				Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
				MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
			},
		},
	}

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, expected)
}

func (s *S) TestListMetadataFormatsWithErrorResponse(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T18:44:25Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="badArgument">'identifier' doesn't match required syntax '(?-xism:^[a-z]+:.*$)'</error>
</OAI-PMH>
	`)

	defer server.Close()
	response, err := client.ListMetadataFormats("x")

	expected := &ListMetadataFormatsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		Error: ResponseError{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "error"},
			Message: "'identifier' doesn't match required syntax '(?-xism:^[a-z]+:.*$)'",
			Code:    "badArgument",
		},
		ResponseDate:    "2016-03-26T18:44:25Z",
		MetadataFormats: []MetadataFormat(nil),
	}

	c.Assert(err, NotNil)
	c.Assert(response, DeepEquals, expected)
}

func (s *S) TestIdentify(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T18:49:42Z</responseDate>
  <request verb="Identify">http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <Identify>
    <repositoryName>ECS EPrints Repository</repositoryName>
    <baseURL>http://eprints.ecs.soton.ac.uk/cgi/oai2</baseURL>
    <protocolVersion>2.0</protocolVersion>
    <adminEmail>cjg@ecs.soton.ac.uk</adminEmail>
    <earliestDatestamp>2011-09-23T08:52:33Z</earliestDatestamp>
    <deletedRecord>persistent</deletedRecord>
    <granularity>YYYY-MM-DDThh:mm:ssZ</granularity>
    <compression>gzip</compression>
    <description>
      <oai-identifier xmlns="http://www.openarchives.org/OAI/2.0/oai-identifier" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai-identifier http://www.openarchives.org/OAI/2.0/oai-identifier.xsd">
	<scheme>oai</scheme>
	<repositoryIdentifier>eprints.ecs.soton.ac.uk</repositoryIdentifier>
	<delimiter>:</delimiter>
	<sampleIdentifier>oai:eprints.ecs.soton.ac.uk:22830</sampleIdentifier>
      </oai-identifier>
    </description>
  </Identify>
</OAI-PMH>
	`)

	defer server.Close()
	response, err := client.Identify()

	expected := &IdentifyResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "Identify",
		},
		Error:        ResponseError{XMLName: xml.Name{Space: "", Local: ""}, Message: "", Code: ""},
		ResponseDate: "2016-03-26T18:49:42Z",
		Identify: Identify{
			XMLName:           xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "Identify"},
			RepositoryName:    "ECS EPrints Repository",
			BaseURL:           "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			ProtocolVersion:   "2.0",
			EarliestDatestamp: "2011-09-23T08:52:33Z",
			DeletedRecord:     "persistent",
			Granularity:       "YYYY-MM-DDThh:mm:ssZ",
			AdminEmail:        "cjg@ecs.soton.ac.uk",
			Compression:       "gzip",
		},
	}

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, expected)
}

func (s *S) TestIdentifyWithErrorResponse(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T19:10:11Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="badArgument">Unrecognised argument 'x'</error>
</OAI-PMH>
	`)

	defer server.Close()
	response, err := client.Identify()

	expected := &IdentifyResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		Error: ResponseError{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "error"},
			Message: "Unrecognised argument 'x'",
			Code:    "badArgument",
		},
		ResponseDate: "2016-03-26T19:10:11Z",
		Identify:     Identify{},
	}

	c.Assert(err, NotNil)
	c.Assert(response, DeepEquals, expected)
}

func (s *S) TestGetRecord(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T19:18:07Z</responseDate>
  <request verb="GetRecord" identifier="oai:eprints.ecs.soton.ac.uk:1" metadataPrefix="oai_dc">http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <GetRecord>
    <record>
    <header>
      <identifier>oai:eprints.ecs.soton.ac.uk:1</identifier>
      <datestamp>2011-09-23T10:22:12Z</datestamp>
      <setSpec>747970653D636F6E666572656E63655F6974656D</setSpec>
      <setSpec>66756C6C746578743D46414C5345</setSpec></header>
    <metadata>
      <oai_dc:dc xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/" xmlns:dc="http://purl.org/dc/elements/1.1/" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<dc:title>A Real Time Neurofuzzy Modelling and State Estimation Scheme</dc:title>
	<dc:creator>Wu, Z.Q.</dc:creator>
	<dc:creator>Harris, C.J.</dc:creator>
	<dc:description>The authors of this paper.</dc:description>
	<dc:publisher>ISO Press</dc:publisher>
	<dc:contributor>Morabito, F.C.</dc:contributor>
	<dc:date>1997</dc:date>
	<dc:type>Conference or Workshop Item</dc:type>
	<dc:type>NonPeerReviewed</dc:type>
	<dc:identifier> Wu, Z. Q. and Harris.</dc:identifier>
	<dc:subject>CS</dc:subject>
	<dc:format>PDF</dc:format>
	<dc:language>en</dc:language>
	<dc:rights>NA</dc:rights>
	<dc:coverage>x</dc:coverage>
	<dc:source>y</dc:source>
	<dc:relation>http://eprints.ecs.soton.ac.uk/1/</dc:relation></oai_dc:dc></metadata></record>
  </GetRecord>
</OAI-PMH>
	`)

	defer server.Close()

	record := new(DublinCoreRecord)
	response, err := client.GetRecord(record, "oai:eprints.ecs.soton.ac.uk:1", "oai_dc")

	expectedResponse := &GetRecordResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "GetRecord",
		},
		Error: ResponseError{XMLName: xml.Name{Space: "", Local: ""}, Message: "", Code: ""},
		Record: Record{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "record"},
			Header: RecordHeader{
				XMLName:    xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "header"},
				Identifier: "oai:eprints.ecs.soton.ac.uk:1",
				Datestamp:  "2011-09-23T10:22:12Z",
				SetSpec: []SetSpec{
					SetSpec{
						XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "setSpec"},
						Set:     "747970653D636F6E666572656E63655F6974656D",
					},
					SetSpec{
						XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "setSpec"},
						Set:     "66756C6C746578743D46414C5345",
					},
				},
				Status: "",
			},
			Metadata: Metadata{
				Raw: []byte(`
      <oai_dc:dc xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/" xmlns:dc="http://purl.org/dc/elements/1.1/" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<dc:title>A Real Time Neurofuzzy Modelling and State Estimation Scheme</dc:title>
	<dc:creator>Wu, Z.Q.</dc:creator>
	<dc:creator>Harris, C.J.</dc:creator>
	<dc:description>The authors of this paper.</dc:description>
	<dc:publisher>ISO Press</dc:publisher>
	<dc:contributor>Morabito, F.C.</dc:contributor>
	<dc:date>1997</dc:date>
	<dc:type>Conference or Workshop Item</dc:type>
	<dc:type>NonPeerReviewed</dc:type>
	<dc:identifier> Wu, Z. Q. and Harris.</dc:identifier>
	<dc:subject>CS</dc:subject>
	<dc:format>PDF</dc:format>
	<dc:language>en</dc:language>
	<dc:rights>NA</dc:rights>
	<dc:coverage>x</dc:coverage>
	<dc:source>y</dc:source>
	<dc:relation>http://eprints.ecs.soton.ac.uk/1/</dc:relation></oai_dc:dc>`),
			},
		},
	}

	expectedRecord := &DublinCoreRecord{
		XMLName:      xml.Name{Space: "http://www.openarchives.org/OAI/2.0/oai_dc/", Local: "dc"},
		Titles:       []string{"A Real Time Neurofuzzy Modelling and State Estimation Scheme"},
		Creators:     []string{"Wu, Z.Q.", "Harris, C.J."},
		Subjects:     []string{"CS"},
		Descriptions: []string{"The authors of this paper."},
		Publishers:   []string{"ISO Press"},
		Contributors: []string{"Morabito, F.C."},
		Dates:        []string{"1997"},
		Types:        []string{"Conference or Workshop Item", "NonPeerReviewed"},
		Formats:      []string{"PDF"},
		Identifiers:  []string{" Wu, Z. Q. and Harris."},
		Sources:      []string{"y"},
		Languages:    []string{"en"},
		Relations:    []string{"http://eprints.ecs.soton.ac.uk/1/"},
		Coverages:    []string{"x"},
		Rights:       []string{"NA"},
	}

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, expectedResponse)
	c.Assert(record, DeepEquals, expectedRecord)
}

func (s *S) TestGetRecordWithErrorResponse(c *C) {
	server, client := mockClient(200, `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-27T17:54:02Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="idDoesNotExist">'oai:eprints.ecs.soton.ac.uk:99999' is not a valid item in this repository</error>
</OAI-PMH>
	`)

	defer server.Close()

	record := new(DublinCoreRecord)
	response, err := client.GetRecord(record, "oai:eprints.ecs.soton.ac.uk:99999", "oai_dc")

	expectedResponse := &GetRecordResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		Error: ResponseError{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "error"},
			Message: "'oai:eprints.ecs.soton.ac.uk:99999' is not a valid item in this repository",
			Code:    "idDoesNotExist",
		},
		Record: Record{},
	}

	expectedRecord := new(DublinCoreRecord)

	c.Assert(err, NotNil)
	c.Assert(response, DeepEquals, expectedResponse)
	c.Assert(record, DeepEquals, expectedRecord)
}
