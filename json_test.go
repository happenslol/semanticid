package semanticid

import (
	"encoding/json"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

type jsonTestSuite struct{}

var _ = Suite(&jsonTestSuite{})

func TestJSON(t *testing.T) { TestingT(t) }

func (s *jsonTestSuite) TestMarshal(c *C) {
	sid, err := NewDefault()
	c.Assert(err, IsNil)

	m, err := json.Marshal(sid)
	c.Assert(err, IsNil)

	mStr := string(m)
	expected := fmt.Sprintf("\"%s\"", sid.String())
	c.Assert(mStr, Equals, expected)

	n, err := json.Marshal(SemanticID{})
	c.Assert(err, IsNil)

	nStr := string(n)
	expected = "null"
	c.Assert(nStr, Equals, expected)
}

func (s *jsonTestSuite) TestUnmarshal(c *C) {
	sid, err := NewDefault()
	c.Assert(err, IsNil)

	jsonStr := fmt.Sprintf("\"%s\"", sid.String())

	var result SemanticID
	err = json.Unmarshal([]byte(jsonStr), &result)
	c.Assert(err, IsNil)

	c.Assert(result.Namespace, Equals, sid.Namespace)
	c.Assert(result.Collection, Equals, sid.Collection)
	c.Assert(result.ID, Equals, sid.ID)

	invalidUUID := "\"namespace:collection:123456789\""
	var invalidUUIDResult SemanticID
	err = json.Unmarshal([]byte(invalidUUID), &invalidUUIDResult)
	c.Assert(err, Not(IsNil))
	c.Assert(invalidUUIDResult.IsNil(), Equals, true)

	invalidSID := "\"123456789\""
	var invalidSIDResult SemanticID
	err = json.Unmarshal([]byte(invalidSID), &invalidSIDResult)
	c.Assert(err, Not(IsNil))
	c.Assert(invalidSIDResult.IsNil(), Equals, true)
}
