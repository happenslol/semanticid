package semanticid

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var rType = reflect.TypeOf(&SemanticID{})

// BSONSemanticIDCodec is a mongodb ValueCodec for
// encoding and decoding SemanticIDs to and from BSON.
type BSONSemanticIDCodec struct{}

var _ bsoncodec.ValueEncoder = &BSONSemanticIDCodec{}
var _ bsoncodec.ValueDecoder = &BSONSemanticIDCodec{}

// EncodeValue implements the ValueEncoder interface.
func (*BSONSemanticIDCodec) EncodeValue(
	ec bsoncodec.EncodeContext,
	vw bsonrw.ValueWriter,
	val reflect.Value,
) error {
	if val.Type() != reflect.TypeOf(&SemanticID{}) {
		return bsoncodec.ValueEncoderError{
			Name:     "SemanticIDEncodeValue",
			Types:    []reflect.Type{rType},
			Received: val,
		}
	}

	return vw.WriteString(val.String())
}

// DecodeValue implements the ValueDecoder interface.
func (*BSONSemanticIDCodec) DecodeValue(
	dc bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) error {
	if !val.CanSet() || val.Type() != rType {
		return bsoncodec.ValueDecoderError{
			Name:     "SemanticIDDecodeValue",
			Types:    []reflect.Type{rType},
			Received: val,
		}
	}

	if vr.Type() != bsontype.String {
		return fmt.Errorf("cannot decode %v into a semanticid", vr.Type())
	}

	str, err := vr.ReadString()
	if err != nil {
		return err
	}

	parsed, err := FromString(str)
	if err != nil {
		return err
	}

	val.FieldByName("Namespace").SetString(parsed.Namespace)
	val.FieldByName("Collection").SetString(parsed.Collection)
	val.FieldByName("UUID").SetString(parsed.UUID)

	return nil
}
