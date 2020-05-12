package semanticid

import (
	"crypto/rand"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/oklog/ulid"
	. "gopkg.in/check.v1"
)

type idProviderTestSuite struct{}

var _ = Suite(&idProviderTestSuite{})

func (s *idProviderTestSuite) SetUpTest(c *C) {}

func TestIDProvider(t *testing.T) { TestingT(t) }

func (s *idProviderTestSuite) TestULID(c *C) {
	validULID := ulid.MustNew(0, rand.Reader).String()
	invalidULID := "1234"

	up := NewULIDProvider()

	gen, err := up.Generate()
	c.Assert(err, IsNil)

	_, err = ulid.Parse(gen)
	c.Assert(err, IsNil)

	valid := up.Validate(validULID)
	c.Assert(valid, IsNil)

	invalid := up.Validate(invalidULID)
	c.Assert(invalid, NotNil)
}

func (s *idProviderTestSuite) TestUUID(c *C) {
	validUUID := uuid.Must(uuid.NewV4()).String()
	invalidUUID := "1234"

	up := NewUUIDProvider()

	gen, err := up.Generate()
	c.Assert(err, IsNil)

	_, err = uuid.FromString(gen)
	c.Assert(err, IsNil)

	valid := up.Validate(validUUID)
	c.Assert(valid, IsNil)

	invalid := up.Validate(invalidUUID)
	c.Assert(invalid, NotNil)
}
