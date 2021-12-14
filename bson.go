package semanticid

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var pointerType = reflect.TypeOf(&SemanticID{})
var rawType = reflect.TypeOf(SemanticID{})

// BSONSemanticIDCodec is a mongodb ValueCodec for
// encoding and decoding SemanticIDs to and from BSON.
type BSONSemanticIDCodec struct{}

// BSONSemanticIDCodec is a mongodb ValueCodec for
// encoding and decoding SemanticIDs to and from BSON.
type BSONSemanticIDPointerCodec struct{}

var _ bsoncodec.ValueEncoder = &BSONSemanticIDCodec{}
var _ bsoncodec.ValueDecoder = &BSONSemanticIDCodec{}

var _ bsoncodec.ValueEncoder = &BSONSemanticIDPointerCodec{}
var _ bsoncodec.ValueDecoder = &BSONSemanticIDPointerCodec{}

// EncodeValue implements the ValueEncoder interface.
func (*BSONSemanticIDCodec) EncodeValue(
	ec bsoncodec.EncodeContext,
	vw bsonrw.ValueWriter,
	val reflect.Value,
) error {
	if val.Type() != rawType {
		return bsoncodec.ValueEncoderError{
			Name:     "SemanticIDEncodeValue",
			Types:    []reflect.Type{rawType},
			Received: val,
		}
	}

	isNilMethod := val.MethodByName("IsNil")
	isNilResult := isNilMethod.Call([]reflect.Value{})[0].Bool()
	if isNilResult {
		return vw.WriteNull()
	}

	strMethod := val.MethodByName("String")
	strResult := strMethod.Call([]reflect.Value{})[0].String()

	return vw.WriteString(strResult)
}

// DecodeValue implements the ValueDecoder interface.
func (*BSONSemanticIDCodec) DecodeValue(
	dc bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) error {
	if !val.CanSet() || val.Type() != rawType {
		return bsoncodec.ValueDecoderError{
			Name:     "SemanticIDDecodeValue",
			Types:    []reflect.Type{rawType},
			Received: val,
		}
	}

	if vr.Type() == bsontype.Null || vr.Type() == bsontype.Undefined {
		val.Set(reflect.ValueOf(SemanticID{}))
		_ = vr.ReadNull()
		return nil
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
	val.FieldByName("ID").SetString(parsed.ID)

	return nil
}

// EncodeValue implements the ValueEncoder interface.
func (*BSONSemanticIDPointerCodec) EncodeValue(
	ec bsoncodec.EncodeContext,
	vw bsonrw.ValueWriter,
	val reflect.Value,
) error {
	if val.IsNil() {
		return vw.WriteNull()
	}

	if val.Type() != pointerType {
		return bsoncodec.ValueEncoderError{
			Name:     "SemanticIDEncodeValue",
			Types:    []reflect.Type{pointerType},
			Received: val,
		}
	}

	isNilMethod := val.MethodByName("IsNil")
	isNilResult := isNilMethod.Call([]reflect.Value{})[0].Bool()
	if isNilResult {
		return vw.WriteString("")
	}

	strMethod := val.MethodByName("String")
	strResult := strMethod.Call([]reflect.Value{})[0].String()

	return vw.WriteString(strResult)
}

// DecodeValue implements the ValueDecoder interface.
func (*BSONSemanticIDPointerCodec) DecodeValue(
	dc bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) error {
	if !val.CanSet() || val.Type() != pointerType {
		return bsoncodec.ValueDecoderError{
			Name:     "SemanticIDDecodeValue",
			Types:    []reflect.Type{pointerType},
			Received: val,
		}
	}

	if vr.Type() == bsontype.Null || vr.Type() == bsontype.Undefined {
		nilValue := reflect.Zero(reflect.TypeOf(&SemanticID{}))
		val.Set(nilValue)
		_ = vr.ReadNull()
		return nil
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

	// NOTE: We need to set this to an empty value, in case
	// the pointer is null. (We're gonna overwrite the existing
	// value anyways, if there is one.)
	val.Set(reflect.ValueOf(&SemanticID{}))
	el := val.Elem()

	el.FieldByName("Namespace").SetString(parsed.Namespace)
	el.FieldByName("Collection").SetString(parsed.Collection)
	el.FieldByName("ID").SetString(parsed.ID)

	return nil
}
