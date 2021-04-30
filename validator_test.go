package semanticid

import (
	"testing"

	"github.com/go-playground/validator/v10"
	. "gopkg.in/check.v1"
)

var _ = Suite(&validatorTestSuite{})
var val *validator.Validate

type TestCollectionValidation struct {
	ID SemanticID `validate:"sidcol=entities"`
}

type TestArrayCollectionValidation struct {
	MultiIDs []SemanticID `validate:"sidcol=entities"`
}

type TestPointerCollectionValidation struct {
	PointerID *SemanticID `validate:"sidcol=entities"`
}

type TestPointerArrayCollectionValidation struct {
	MultiPointerIDs []SemanticID `validate:"sidcol=entities"`
}

type validatorTestSuite struct{}

func TestValidator(t *testing.T) { TestingT(t) }

func (s *validatorTestSuite) SetUpSuite(c *C) {
	val = validator.New()
	RegisterValidation(val)
}

func (s *validatorTestSuite) TestValidator(c *C) {
	sid, err := NewWithCollection("entities")
	c.Assert(err, IsNil)

	wrongCollectionSID, err := NewWithCollection("wrong")
	c.Assert(err, IsNil)

	collectionNil := TestCollectionValidation{}
	c.Assert(val.Struct(collectionNil), NotNil)

	collectionInvalid := TestCollectionValidation{ID: wrongCollectionSID}
	c.Assert(val.Struct(collectionInvalid), NotNil)

	collectionValid := TestCollectionValidation{ID: sid}
	c.Assert(val.Struct(collectionValid), IsNil)
}
