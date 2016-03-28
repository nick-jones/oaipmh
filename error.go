package oaipmh

import (
	"fmt"
)

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

func (e Error) Empty() bool {
	return e.Code == "" && e.Message == ""
}
