package oaipmh

import (
	"encoding/xml"
	. "gopkg.in/check.v1"
)

type errorSuite struct{}

var _ = Suite(&errorSuite{})

func (s *errorSuite) TestErrorReturnsFormattedString(c *C) {
	err := &Error{
		XMLName: xml.Name{},
		Message: "foo",
		Code:    "c",
	}

	c.Assert(err.Error(), Equals, "c: foo")
}

func (s *errorSuite) TestErrorIndicatesEmptyWhenNoCodeOrMessageAreAssigned(c *C) {
	err := &Error{
		XMLName: xml.Name{},
		Message: "",
		Code:    "",
	}

	c.Assert(err.Empty(), Equals, true)
}

func (s *errorSuite) TestErrorIndicatesNonEmptyWhenCodeIsAssigned(c *C) {
	err := &Error{
		XMLName: xml.Name{},
		Message: "",
		Code:    "x",
	}

	c.Assert(err.Empty(), Equals, false)
}

func (s *errorSuite) TestErrorIndicatesNonEmptyWhenMessageIsAssigned(c *C) {
	err := &Error{
		XMLName: xml.Name{},
		Message: "x",
		Code:    "",
	}

	c.Assert(err.Empty(), Equals, false)
}
