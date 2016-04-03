package oaipmh

import (
	"encoding/xml"
	"fmt"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

type clientSuite struct{}

var _ = Suite(&clientSuite{})

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

func (s *clientSuite) TestListMetadataFormats(c *C) {
	raw := `
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
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	request := &ListMetadataFormatsOptions{}
	metadataFormats, _, err := client.ListMetadataFormats(request)

	expectedMetadataFormats := &ListMetadataFormatsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "ListMetadataFormats",
		},
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
	c.Assert(metadataFormats, DeepEquals, expectedMetadataFormats)
}

func (s *clientSuite) TestListMetadataFormatsWithErrorResponse(c *C) {
	raw := `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T18:44:25Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="badArgument">'identifier' doesn't match required syntax '(?-xism:^[a-z]+:.*$)'</error>
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	request := &ListMetadataFormatsOptions{"x"}
	metadataFormats, _, err := client.ListMetadataFormats(request)

	expectedMetadataFormats := &ListMetadataFormatsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		ResponseDate:    "2016-03-26T18:44:25Z",
		MetadataFormats: []MetadataFormat(nil),
	}

	c.Assert(err, NotNil)
	c.Assert(metadataFormats, DeepEquals, expectedMetadataFormats)
}

func (s *clientSuite) TestIdentify(c *C) {
	raw := `
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
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	identity, _, err := client.Identify()

	expectedIdentity := &IdentifyResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "Identify",
		},
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
	c.Assert(identity, DeepEquals, expectedIdentity)
}

func (s *clientSuite) TestIdentifyWithErrorResponse(c *C) {
	raw := `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-26T19:10:11Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="badArgument">Unrecognised argument 'x'</error>
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	identity, _, err := client.Identify()

	expectedIdentity := &IdentifyResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		ResponseDate: "2016-03-26T19:10:11Z",
		Identify:     Identify{},
	}

	c.Assert(err, NotNil)
	c.Assert(identity, DeepEquals, expectedIdentity)
}

func (s *clientSuite) TestGetRecord(c *C) {
	raw := `
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
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	metadata := new(DublinCoreRecord)
	request := &GetRecordOptions{"oai:eprints.ecs.soton.ac.uk:1", "oai_dc"}
	record, _, err := client.GetRecord(request, metadata)

	expectedRecord := &GetRecordResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "GetRecord",
		},
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

	expectedMetadata := &DublinCoreRecord{
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
	c.Assert(record, DeepEquals, expectedRecord)
	c.Assert(metadata, DeepEquals, expectedMetadata)
}

func (s *clientSuite) TestGetRecordWithErrorResponse(c *C) {
	raw := `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-27T17:54:02Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="idDoesNotExist">'oai:eprints.ecs.soton.ac.uk:99999' is not a valid item in this repository</error>
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	metadata := new(DublinCoreRecord)
	request := &GetRecordOptions{"oai:eprints.ecs.soton.ac.uk:99999", "oai_dc"}
	record, _, err := client.GetRecord(request, metadata)

	expectedRecord := &GetRecordResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		Record: Record{},
	}

	expectedMetadata := &DublinCoreRecord{}

	c.Assert(err, NotNil)
	c.Assert(record, DeepEquals, expectedRecord)
	c.Assert(metadata, DeepEquals, expectedMetadata)
}

func (s *clientSuite) TestListRecords(c *C) {
	raw := `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-27T18:20:04Z</responseDate>
  <request verb="ListRecords" metadataPrefix="oai_dc">http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <ListRecords>
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
        <dc:description>The authors of this paper</dc:description>
        <dc:publisher>ISO Press</dc:publisher>
        <dc:contributor>Morabito, F.C.</dc:contributor>
        <dc:date>1997</dc:date>
        <dc:type>Conference or Workshop Item</dc:type>
        <dc:type>NonPeerReviewed</dc:type>
        <dc:identifier>Wu, Z. Q. and Harris</dc:identifier>
        <dc:relation>http://eprints.ecs.soton.ac.uk/1/</dc:relation></oai_dc:dc></metadata></record>
    <resumptionToken expirationDate="2016-03-28T18:20:04Z">metadataPrefix%3Doai_dc%26offset%3D101</resumptionToken>
  </ListRecords>
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	metadatas := new(DublinCoreRecords)
	options := &ListOptions{"oai_dc", time.Time{}, time.Time{}, "", ""}
	records, _, err := client.ListRecords(options, metadatas)

	expectedRecords := &ListRecordsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "ListRecords",
		},
		Records: []Record{
			Record{
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
        <dc:description>The authors of this paper</dc:description>
        <dc:publisher>ISO Press</dc:publisher>
        <dc:contributor>Morabito, F.C.</dc:contributor>
        <dc:date>1997</dc:date>
        <dc:type>Conference or Workshop Item</dc:type>
        <dc:type>NonPeerReviewed</dc:type>
        <dc:identifier>Wu, Z. Q. and Harris</dc:identifier>
        <dc:relation>http://eprints.ecs.soton.ac.uk/1/</dc:relation></oai_dc:dc>`),
				},
			},
		},
		ResumptionToken: ResumptionToken{
			XMLName:        xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "resumptionToken"},
			ExpirationDate: "2016-03-28T18:20:04Z",
			Value:          "metadataPrefix%3Doai_dc%26offset%3D101",
		},
	}

	expectedMetadatas := &DublinCoreRecords{
		Records: []DublinCoreRecord{
			DublinCoreRecord{
				XMLName:      xml.Name{Space: "http://www.openarchives.org/OAI/2.0/oai_dc/", Local: "dc"},
				Titles:       []string{"A Real Time Neurofuzzy Modelling and State Estimation Scheme"},
				Creators:     []string{"Wu, Z.Q.", "Harris, C.J."},
				Subjects:     []string(nil),
				Descriptions: []string{"The authors of this paper"},
				Publishers:   []string{"ISO Press"},
				Contributors: []string{"Morabito, F.C."},
				Dates:        []string{"1997"},
				Types:        []string{"Conference or Workshop Item", "NonPeerReviewed"},
				Formats:      []string(nil),
				Identifiers:  []string{"Wu, Z. Q. and Harris"},
				Sources:      []string(nil),
				Languages:    []string(nil),
				Relations:    []string{"http://eprints.ecs.soton.ac.uk/1/"},
				Coverages:    []string(nil),
				Rights:       []string(nil),
			},
		},
	}

	c.Assert(err, IsNil)
	c.Assert(records, DeepEquals, expectedRecords)
	c.Assert(metadatas, DeepEquals, expectedMetadatas)
}

func (s *clientSuite) TestListRecordsWithErrorResponse(c *C) {
	raw := `
<?xml version='1.0' encoding='UTF-8'?>
<?xml-stylesheet type='text/xsl' href='/oai2.xsl' ?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
  <responseDate>2016-03-27T18:37:10Z</responseDate>
  <request>http://eprints.ecs.soton.ac.uk/cgi/oai2</request>
  <error code="cannotDisseminateFormat">Record not available as metadata type: 'oai_dd'</error>
</OAI-PMH>`

	server, client := mockClient(200, raw)

	defer server.Close()
	metadatas := new(DublinCoreRecords)
	request := &ListOptions{"oai_dd", time.Time{}, time.Time{}, "", ""}
	records, _, err := client.ListRecords(request, metadatas)

	expectedRecords := &ListRecordsResponse{
		XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "OAI-PMH"},
		InterpretedRequest: InterpretedRequest{
			XMLName: xml.Name{Space: "http://www.openarchives.org/OAI/2.0/", Local: "request"},
			BaseURL: "http://eprints.ecs.soton.ac.uk/cgi/oai2",
			Verb:    "",
		},
		Records:         []Record(nil),
		ResumptionToken: ResumptionToken{},
	}

	expectedMetadatas := &DublinCoreRecords{
		Records: []DublinCoreRecord(nil),
	}

	c.Assert(err, NotNil)
	c.Assert(records, DeepEquals, expectedRecords)
	c.Assert(metadatas, DeepEquals, expectedMetadatas)
}
