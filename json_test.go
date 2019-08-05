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
	c.Assert(result.UUID, Equals, sid.UUID)
}
