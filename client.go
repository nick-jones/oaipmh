package oaipmh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const formatISO8601 string = "%04d-%02d-%02dT%02d:%02d:%02dZ"

type HTTPResponse struct {
	StatusCode int
	Raw        []byte
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}, nil
}

func (c *Client) ListMetadataFormats(request *ListMetadataFormatsOptions) (*ListMetadataFormatsResponse, *HTTPResponse, error) {
	params := prepareParameters("ListMetadataFormats", map[string]string{
		"identifier": request.Identifier,
	})

	response := new(ListMetadataFormatsResponse)
	httpResponse, err := c.fetchXML(params, response)

	return response, httpResponse, err
}

func (c *Client) Identify() (*IdentifyResponse, *HTTPResponse, error) {
	params := prepareParameters("Identify", map[string]string{})
	response := new(IdentifyResponse)
	httpResponse, err := c.fetchXML(params, response)

	return response, httpResponse, err
}

func (c *Client) GetRecord(request *GetRecordOptions, record interface{}) (*GetRecordResponse, *HTTPResponse, error) {
	params := prepareParameters("GetRecord", map[string]string{
		"identifier":     request.Identifier,
		"metadataPrefix": request.MetadataPrefix,
	})

	response := new(GetRecordResponse)
	httpResponse, err := c.fetchXML(params, response)

	if err != nil {
		return response, httpResponse, err
	}

	return response, httpResponse, unmarshalRecord(response.Record, record)
}

func (c *Client) ListRecords(request *ListRecordsOptions, records interface{}) (*ListRecordsResponse, *HTTPResponse, error) {
	params := prepareParameters("ListRecords", map[string]string{
		"metadataPrefix":  request.MetadataPrefix,
		"from":            formatDateTime(request.From),
		"until":           formatDateTime(request.Until),
		"set":             request.Set,
		"resumptionToken": request.ResumptionToken,
	})

	response := new(ListRecordsResponse)
	httpResponse, err := c.fetchXML(params, response)

	if err != nil {
		return response, httpResponse, err
	}

	return response, httpResponse, unmarshalRecords(response.Records, records)
}

func (c *Client) fetch(params url.Values) (*HTTPResponse, error) {
	query := params.Encode()
	path := fmt.Sprintf("%s?%s", c.baseURL, query)
	res, err := c.http.Get(path)

	if err != nil {
		return &HTTPResponse{}, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	httpResponse := &HTTPResponse{res.StatusCode, contents}

	if httpResponse.StatusCode >= 400 {
		err = errors.New("Unsuccessful request")
	}

	return httpResponse, err
}

func (c *Client) fetchXML(params url.Values, into interface{}) (*HTTPResponse, error) {
	httpResponse, err := c.fetch(params)

	if err != nil {
		return httpResponse, err
	}

	return httpResponse, unmarshalResponse(httpResponse.Raw, into)
}

func prepareParameters(verb string, options map[string]string) url.Values {
	params := url.Values{}
	params.Add("verb", verb)

	for key, value := range options {
		if value != "" {
			params.Add(key, value)
		}
	}

	return params
}

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return fmt.Sprintf(formatISO8601, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
