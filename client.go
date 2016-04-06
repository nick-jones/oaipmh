package oaipmh

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const formatISO8601 string = "%04d-%02d-%02dT%02d:%02d:%02dZ"

type Client struct {
	baseURL string
	http    *http.Client
}

type HTTPResponse struct {
	StatusCode int
	Raw        []byte
}

func NewClient(baseURL string) (*Client, error) {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}, nil
}

func (c *Client) ListMetadataFormats(options *ListMetadataFormatsOptions) (*ListMetadataFormatsResponse, *HTTPResponse, error) {
	params := prepareParameters("ListMetadataFormats", map[string]string{
		"identifier": options.Identifier,
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

func (c *Client) GetRecord(options *GetRecordOptions, record interface{}) (*GetRecordResponse, *HTTPResponse, error) {
	params := prepareParameters("GetRecord", map[string]string{
		"identifier":     options.Identifier,
		"metadataPrefix": options.MetadataPrefix,
	})

	response := new(GetRecordResponse)
	httpResponse, err := c.fetchXML(params, response)

	if err != nil {
		return response, httpResponse, err
	}

	return response, httpResponse, unmarshalRecord(response.Record, record)
}

func (c *Client) ListRecords(options *ListOptions, records interface{}) (*ListRecordsResponse, *HTTPResponse, error) {
	params := prepareParameters("ListRecords", map[string]string{
		"metadataPrefix":  options.MetadataPrefix,
		"from":            formatDateTime(options.From),
		"until":           formatDateTime(options.Until),
		"set":             options.Set,
		"resumptionToken": options.ResumptionToken,
	})

	response := new(ListRecordsResponse)
	httpResponse, err := c.fetchXML(params, response)

	if err != nil {
		return response, httpResponse, err
	}

	return response, httpResponse, unmarshalRecords(response.Records, records)
}

func (c *Client) ListIdentifiers(options *ListOptions) (*ListIdentifiersResponse, *HTTPResponse, error) {
	params := prepareParameters("ListIdentifiers", map[string]string{
		"metadataPrefix":  options.MetadataPrefix,
		"from":            formatDateTime(options.From),
		"until":           formatDateTime(options.Until),
		"set":             options.Set,
		"resumptionToken": options.ResumptionToken,
	})

	response := new(ListIdentifiersResponse)
	httpResponse, err := c.fetchXML(params, response)

	return response, httpResponse, err
}

func (c *Client) ListSets(options *ListSetsOptions) (*ListSetsResponse, *HTTPResponse, error) {
	params := prepareParameters("ListSets", map[string]string{
		"resumptionToken": options.ResumptionToken,
	})

	response := new(ListSetsResponse)
	httpResponse, err := c.fetchXML(params, response)

	return response, httpResponse, err
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

func unmarshalResponse(data []byte, into interface{}) error {
	if err := xml.Unmarshal(data, into); err != nil {
		return err
	}

	responseError := new(ResponseError)

	if err := xml.Unmarshal(data, responseError); err != nil {
		return err
	}

	if !responseError.Error.Empty() {
		return responseError.Error
	}

	return nil
}

func unmarshalRecord(record Record, into interface{}) error {
	typ := reflect.TypeOf(into).Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("Non-struct provided")
	}

	return xml.Unmarshal(record.Metadata.Raw, into)
}

func unmarshalRecords(records []Record, into interface{}) error {
	pointer := reflect.ValueOf(into)
	elem := pointer.Elem()

	if elem.Kind() != reflect.Struct {
		return errors.New("Non-struct provided")
	}

	field := elem.FieldByName("Records")

	if !field.IsValid() {
		return errors.New("Struct provided must contain `Records` field")
	}

	typ := field.Type().Elem()
	size := len(records)
	slice := reflect.MakeSlice(reflect.SliceOf(typ), size, size)

	for i, item := range records {
		value := reflect.New(typ)
		xml.Unmarshal(item.Metadata.Raw, value.Interface())
		slice.Index(i).Set(value.Elem())
	}

	field.Set(slice)

	return nil
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
