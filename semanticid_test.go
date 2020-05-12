package semanticid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/oklog/ulid"
	. "gopkg.in/check.v1"
)

var _ IDProvider = &TestProvider{}
var testProvider = &TestProvider{}

func (tp *TestProvider) Generate() (string, error) {
	return "1234", nil
}

func (tp *TestProvider) Validate(id string) error {
	return nil
}

type semanticidTestSuite struct{}

var _ = Suite(&semanticidTestSuite{})

func (s *semanticidTestSuite) SetUpTest(c *C) {
	// Reset global defaults
	DefaultNamespace = "namespace"
	DefaultCollection = "collection"
	DefaultIDProvider = NewULIDProvider()
	Separator = "."
}

// Hook up gocheck into the "go test" runner.
func TestSemanticID(t *testing.T) { TestingT(t) }

func (s *semanticidTestSuite) TestNew(c *C) {
	sid, err := NewDefault()

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	sid, err = NewWithNamespace("test")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	sid, err = NewWithCollection("test")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "test")
	c.Assert(sid.ID, Not(HasLen), 0)

	sid, err = New("test-namespace", "test-collection")

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test-namespace")
	c.Assert(sid.Collection, Equals, "test-collection")
	c.Assert(sid.ID, Not(HasLen), 0)
}

func (s *semanticidTestSuite) TestDefault(c *C) {
	DefaultNamespace = "test-default-namespace"
	DefaultCollection = "test-default-collection"
	DefaultIDProvider = testProvider

	sid, err := NewDefault()

	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "test-default-namespace")
	c.Assert(sid.Collection, Equals, "test-default-collection")
	c.Assert(sid.ID, Equals, "1234")
}

func (s *semanticidTestSuite) TestSeparator(c *C) {
	Separator = ":"
	sid, err := NewDefault()
	sidStr := sid.String()

	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(sidStr, "namespace:collection:"), Equals, true)

	Separator = "."

	_, err = NewWithNamespace("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errPartContainsSeparator)

	_, err = NewWithCollection("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errPartContainsSeparator)
}

func (s *semanticidTestSuite) TestFromString(c *C) {
	validID := ulid.MustNew(0, rand.Reader).String()
	invalidID := "0123456789"

	sid, err := FromString(fmt.Sprintf("a.b.%s", validID))
	c.Assert(err, IsNil)
	c.Assert(sid.Namespace, Equals, "a")
	c.Assert(sid.Collection, Equals, "b")
	c.Assert(sid.ID, Equals, validID)

	sid, err = FromString(invalidID)
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errInvalidSID)

	sid, err = FromString(fmt.Sprintf("a.b.%s", invalidID))
	c.Assert(err, Not(IsNil))
	c.Assert(err.(semanticIDError).errCode, Equals, errInvalidID)
}

func (s *semanticidTestSuite) TestString(c *C) {
	sid, err := NewDefault()
	sidStr := sid.String()

	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(sidStr, "namespace.collection."), Equals, true)
}

func (s *semanticidTestSuite) TestIsNil(c *C) {
	sid, err := NewWithNamespace("test.test")
	c.Assert(err, Not(IsNil))
	c.Assert(sid.IsNil(), Equals, true)

	parsed, err := FromString("0123456789")
	c.Assert(err, Not(IsNil))
	c.Assert(parsed.IsNil(), Equals, true)
}

func (s *semanticidTestSuite) TestMust(c *C) {
	sid := Must(NewDefault())
	c.Assert(sid.Namespace, Equals, "namespace")
	c.Assert(sid.Collection, Equals, "collection")
	c.Assert(sid.ID, Not(HasLen), 0)

	defer func() {
		c.Assert(recover(), NotNil)
	}()

	Must(func() (SemanticID, error) {
		return SemanticID{}, errors.New("")
	}())
}

type TestModel struct {
	ID         SemanticID `sid:"testmodels"`
	CustomID   string     `sid:"custom"`
	InvalidTag string     `sid:"test.models"`
}

func (s *semanticidTestSuite) TestModel(c *C) {
	collection, err := CollectionForModel(TestModel{})
	c.Assert(err, IsNil)
	c.Assert(collection, Equals, "testmodels")

	collection, err = CollectionForModel(&TestModel{})
	c.Assert(err, IsNil)
	c.Assert(collection, Equals, "testmodels")

	collection, err = CollectionForModelField(TestModel{}, "CustomID")
	c.Assert(err, IsNil)
	c.Assert(collection, Equals, "custom")

	_, err = CollectionForModelField(TestModel{}, "InvalidTag")
	c.Assert(err, NotNil)

	_, err = CollectionForModelField(TestModel{}, "NonExistant")
	c.Assert(err, NotNil)

	sid, err := NewForModel(TestModel{})
	c.Assert(err, IsNil)
	c.Assert(sid.Collection, Equals, "testmodels")
}
