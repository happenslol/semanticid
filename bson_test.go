package semanticid

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	. "gopkg.in/check.v1"
)

type bsonTestSuite struct{}

var _ = Suite(&bsonTestSuite{})
var reg *bsoncodec.Registry

func TestBSON(t *testing.T) { TestingT(t) }

func (s *bsonTestSuite) SetUpSuite(c *C) {
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)

	rb.RegisterCodec(
		reflect.TypeOf(&SemanticID{}),
		&BSONSemanticIDPointerCodec{},
	)

	rb.RegisterCodec(
		reflect.TypeOf(SemanticID{}),
		&BSONSemanticIDCodec{},
	)

	reg = rb.Build()
}

func (s *bsonTestSuite) TestValue(c *C) {
	sid, err := NewDefault()
	c.Assert(err, IsNil)

	val := bson.M{"id": sid}

	m, err := bson.MarshalWithRegistry(reg, val)
	c.Assert(err, IsNil)

	var result map[string]SemanticID
	err = bson.UnmarshalWithRegistry(reg, m, &result)
	c.Assert(err, IsNil)

	c.Assert(result["id"].Namespace, Equals, sid.Namespace)
	c.Assert(result["id"].Collection, Equals, sid.Collection)
	c.Assert(result["id"].UUID, Equals, sid.UUID)
}

func (s *bsonTestSuite) TestPointer(c *C) {
	sid, err := NewDefault()
	c.Assert(err, IsNil)

	val := bson.M{"id": &sid}

	m, err := bson.MarshalWithRegistry(reg, val)
	c.Assert(err, IsNil)

	var result map[string]*SemanticID
	err = bson.UnmarshalWithRegistry(reg, m, &result)
	c.Assert(err, IsNil)

	c.Assert(result["id"].Namespace, Equals, sid.Namespace)
	c.Assert(result["id"].Collection, Equals, sid.Collection)
	c.Assert(result["id"].UUID, Equals, sid.UUID)
}

func (s *bsonTestSuite) TestNull(c *C) {
	sid := SemanticID{}
	val := bson.M{"id": sid}

	m, err := bson.MarshalWithRegistry(reg, &val)
	c.Assert(err, IsNil)

	var result map[string]SemanticID
	err = bson.UnmarshalWithRegistry(reg, m, &result)

	c.Assert(err, IsNil)
	c.Assert(result["id"].IsNil(), Equals, true)
}
