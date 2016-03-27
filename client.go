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

func (c *Client) ListMetadataFormats(request *ListMetadataFormatsOptions) (*ListMetadataFormatsResponse, error) {
	params := prepareUrlValues("ListMetadataFormats", map[string]string{
		"identifier": request.Identifier,
	})

	response := new(ListMetadataFormatsResponse)

	if err := c.fetchXML(params, response); err != nil {
		return response, err
	}

	return response, errorFromResponse(response.Error)
}

func (c *Client) Identify() (*IdentifyResponse, error) {
	params := prepareUrlValues("Identify", map[string]string{})
	response := new(IdentifyResponse)

	if err := c.fetchXML(params, response); err != nil {
		return response, err
	}

	return response, errorFromResponse(response.Error)
}

func (c *Client) GetRecord(record interface{}, request *GetRecordOptions) (*GetRecordResponse, error) {
	params := prepareUrlValues("GetRecord", map[string]string{
		"identifier":     request.Identifier,
		"metadataPrefix": request.MetadataPrefix,
	})

	response := new(GetRecordResponse)
	data, err := c.fetch(params)

	if err != nil {
		return response, err
	}

	if err = xml.Unmarshal(data, response); err != nil {
		return response, err
	}

	if err = errorFromResponse(response.Error); err != nil {
		return response, err
	}

	return response, unmarshalRecord(response.Record, record)
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

func unmarshalRecord(record Record, into interface{}) error {
	typ := reflect.TypeOf(into).Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("Non-struct provided")
	}

	return xml.Unmarshal(record.Metadata.Raw, into)
}

func prepareUrlValues(verb string, options map[string]string) url.Values {
	params := url.Values{}
	params.Add("verb", verb)

	for key, value := range options {
		if value != "" {
			params.Add(key, value)
		}
	}

	return params
}

func errorFromResponse(responseError ResponseError) error {
	if responseError.Code != "" || responseError.Message != "" {
		return fmt.Errorf("API error: %s (%s)", responseError.Message, responseError.Code)
	}

	return nil
}