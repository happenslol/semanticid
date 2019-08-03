package semanticid

import (
	"testing"
	. "gopkg.in/check.v1"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) SetUpTest(c *C) {
	DefaultNamespace = "namespace"
	DefaultCollection = "collection"
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

func (s *testSuite) TestSeparator(c *C) {}