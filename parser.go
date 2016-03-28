package oaipmh

import (
	"encoding/xml"
	"errors"
	"reflect"
)

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
