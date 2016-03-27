package oaipmh

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

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

func (c *Client) ListMetadataFormats(identifier string) (*ListMetadataFormatsResponse, error) {
	options := map[string]string{
		"identifier": identifier,
	}

	params := c.prepareUrlValues("ListMetadataFormats", options)
	response := new(ListMetadataFormatsResponse)

	if err := c.fetchXML(params, response); err != nil {
		return response, err
	}

	return response, c.errorFromResponse(response.Error)
}

func (c *Client) Identify() (*IdentifyResponse, error) {
	params := c.prepareUrlValues("Identify", map[string]string{})
	response := new(IdentifyResponse)

	if err := c.fetchXML(params, response); err != nil {
		return response, err
	}

	return response, c.errorFromResponse(response.Error)
}

func (c *Client) GetRecord(record interface{}, identifier string, metadataPrefix string) (*GetRecordResponse, error) {
	options := map[string]string{
		"identifier":     identifier,
		"metadataPrefix": metadataPrefix,
	}

	params := c.prepareUrlValues("GetRecord", options)
	response := new(GetRecordResponse)
	data, err := c.fetch(params)

	if err != nil {
		return response, err
	}

	if err = xml.Unmarshal(data, response); err != nil {
		return response, err
	}

	if err = c.unmarshalRecord(response.Record, record); err != nil {
		return response, err
	}

	return response, c.errorFromResponse(response.Error)
}

func (c *Client) unmarshalRecord(record Record, into interface{}) error {
	typ := reflect.TypeOf(into).Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("Non-struct provided")
	}

	return xml.Unmarshal(record.Metadata.Raw, into)
}

func (c *Client) fetch(params url.Values) ([]byte, error) {
	query := params.Encode()
	path := fmt.Sprintf("%s?%s", c.baseURL, query)
	res, err := c.http.Get(path)
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		err = errors.New("Unsuccessful request")
	}

	return contents, err
}

func (c *Client) fetchXML(params url.Values, response interface{}) error {
	res, err := c.fetch(params)

	if err != nil {
		return err
	}

	return xml.Unmarshal(res, response)
}

func (c *Client) prepareUrlValues(verb string, options map[string]string) url.Values {
	params := url.Values{}
	params.Add("verb", verb)

	for key, value := range options {
		if value != "" {
			params.Add(key, value)
		}
	}

	return params
}

func (c *Client) errorFromResponse(responseError ResponseError) error {
	if responseError.Code != "" || responseError.Message != "" {
		return errors.New("Error indicated by endpoint")
	}

	return nil
}
