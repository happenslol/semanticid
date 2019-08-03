package semanticid

import (
	"errors"
	"fmt"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) SetUpTest(c *C) {
	// Reset global defaults
	DefaultNamespace = "namespace"
	DefaultCollection = "collection"
	Separator = "."
}

// Hook up gocheck into the "go test" runner.
func TestSemanticID(t *testing.T) { TestingT(t) }

func (s *testSuite) TestNew(c *C) {
	sid, err := NewDefault()

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.UUID, Not(HasLen), 0)

	sid, err = NewWithNamespace("test")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.UUID, Not(HasLen), 0)

	sid, err = NewWithCollection("test")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "test")
	c.Assert(sid.UUID, Not(HasLen), 0)

	sid, err = New("test-namespace", "test-collection")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test-namespace")
	c.Assert(sid.Collection, Equals, "test-collection")
	c.Assert(sid.UUID, Not(HasLen), 0)
}

func (s *testSuite) TestDefault(c *C) {
	DefaultNamespace = "test-default-namespace"
	DefaultCollection = "test-default-collection"

	sid, err := NewDefault()

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test-default-namespace")
	c.Assert(sid.Collection, Equals, "test-default-collection")
	c.Assert(sid.UUID, Not(HasLen), 0)
}

func (s *testSuite) TestSeparator(c *C) {
	Separator = ":"
	sid, err := NewDefault()
	sidStr := sid.String()

	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(sidStr, "namespace:collection:"), Equals, true)

	Separator = "."

	_, err = NewWithNamespace("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errPartContainsSeparator)
	c.Assert(strings.Contains(err.Error(), "namespace"), Equals, true)

	_, err = NewWithCollection("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errPartContainsSeparator)
	c.Assert(strings.Contains(err.Error(), "collection"), Equals, true)
}

func (s *testSuite) TestFromString(c *C) {
	validUUID := "fd16ef44-3187-4b0d-8f18-93e7f7a5d88a"
	invalidUUID := "0123456789"

	sid, err := FromString(fmt.Sprintf("a.b.%s", validUUID))
	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "a")
	c.Assert(sid.Collection, Equals, "b")
	c.Assert(sid.UUID, Equals, validUUID)

	sid, err = FromString(invalidUUID)
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errInvalidSID)

	sid, err = FromString(fmt.Sprintf("a.b.%s", invalidUUID))
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errInvalidUUID)
}

func (s *testSuite) TestString(c *C) {
	sid, err := NewDefault()
	sidStr := sid.String()

	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(sidStr, "namespace.collection."), Equals, true)
}

func (s *testSuite) TestIsNil(c *C) {
	sid, err := NewWithNamespace("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(sid.IsNil(), Equals, true)

	parsed, err := FromString("0123456789")
	c.Assert(err, Not(IsNil))
	c.Assert(parsed.IsNil(), Equals, true)
}

func (s *testSuite) TestMust(c *C) {
	sid := Must(NewDefault())
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.UUID, Not(HasLen), 0)

	defer func() {
		c.Assert(recover(), NotNil)
	}()

	Must(func() (SemanticID, error) {
		return SemanticID{}, errors.New("")
	}())
}
