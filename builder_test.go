package semanticid

import (
	"testing"

	. "gopkg.in/check.v1"
)

type builderTestSuite struct{}

var _ = Suite(&builderTestSuite{})

func (s *builderTestSuite) SetUpTest(c *C) {
	// Reset global defaults
	DefaultNamespace = "namespace"
	DefaultCollection = "collection"
	DefaultIDProvider = NewULIDProvider()
	Separator = "."
}

func TestBuilder(t *testing.T) { TestingT(t) }

type TestProvider struct{}

func (s *builderTestSuite) TestBuild(c *C) {
	invalidSID := "a.b.1234"
	validSID := Must(NewDefault()).String()

	sid, err := Builder().Build()
	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	sid, err = Builder().WithIDProvider(testProvider).Build()
	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Equals, "1234")

	sid, err = Builder().
		WithNamespace("test-namespace").
		WithCollection("test-collection").
		Build()

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test-namespace")
	c.Assert(sid.Collection, Equals, "test-collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	sid, err = Builder().FromString(validSID).Build()
	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	_, err = Builder().FromString(invalidSID).Build()
	c.Assert(err, NotNil)

	_, err = Builder().FromString(invalidSID).NoValidate().Build()
	c.Assert(err, IsNil)
}
