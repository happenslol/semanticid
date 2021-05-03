package semanticid_test

import (
	"reflect"

	"github.com/happenslol/semanticid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

var _ = Describe("bson", func() {
	var (
		reg            *bsoncodec.Registry
		sidValue       semanticid.SemanticID
		sidPointer     *semanticid.SemanticID
		zeroSID        semanticid.SemanticID
		zeroSIDPointer *semanticid.SemanticID
	)

	BeforeEach(func() {
		rb := bsoncodec.NewRegistryBuilder()
		bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
		bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)

		rb.RegisterCodec(
			reflect.TypeOf(&semanticid.SemanticID{}),
			&semanticid.BSONSemanticIDPointerCodec{},
		)

		rb.RegisterCodec(
			reflect.TypeOf(semanticid.SemanticID{}),
			&semanticid.BSONSemanticIDCodec{},
		)

		reg = rb.Build()

		var err error
		sidValue, err = semanticid.NewDefault()
		Expect(err).To(BeNil())

		sidPointerValue, err := semanticid.NewDefault()
		Expect(err).To(BeNil())

		sidPointer = &sidPointerValue

		zeroSID = semanticid.SemanticID{}
		zeroSIDPointer = &semanticid.SemanticID{}
	})

	Describe("Marshalling and unmarshalling semanticids as bson", func() {
		Context("with a semanticid value", func() {
			It("should work correctly", func() {
				val := bson.M{"id": sidValue}
				m, err := bson.MarshalWithRegistry(reg, val)
				Expect(err).To(BeNil())

				var result map[string]semanticid.SemanticID
				err = bson.UnmarshalWithRegistry(reg, m, &result)
				Expect(err).To(BeNil())

				Expect(result["id"].Namespace).To(Equal(sidValue.Namespace))
				Expect(result["id"].Collection).To(Equal(sidValue.Collection))
				Expect(result["id"].ID).To(Equal(sidValue.ID))
			})
		})

		Context("with a semanticid pointer", func() {
			It("should work correctly", func() {
				val := bson.M{"id": sidPointer}
				m, err := bson.MarshalWithRegistry(reg, val)
				Expect(err).To(BeNil())

				var result map[string]*semanticid.SemanticID
				err = bson.UnmarshalWithRegistry(reg, m, &result)
				Expect(err).To(BeNil())

				Expect(result["id"].Namespace).To(Equal(sidPointer.Namespace))
				Expect(result["id"].Collection).To(Equal(sidPointer.Collection))
				Expect(result["id"].ID).To(Equal(sidPointer.ID))
			})
		})

		Context("with a zero semanticid value", func() {
			It("should work correctly", func() {
				val := bson.M{"id": zeroSID}
				m, err := bson.MarshalWithRegistry(reg, val)
				Expect(err).To(BeNil())

				var result map[string]semanticid.SemanticID
				err = bson.UnmarshalWithRegistry(reg, m, &result)
				Expect(err).To(BeNil())

				Expect(result["id"].IsNil()).To(BeTrue())
			})
		})

		Context("with a zero semanticid pointer", func() {
			It("should work correctly", func() {
				val := bson.M{"id": zeroSIDPointer}
				m, err := bson.MarshalWithRegistry(reg, val)
				Expect(err).To(BeNil())

				var result map[string]*semanticid.SemanticID
				err = bson.UnmarshalWithRegistry(reg, m, &result)
				Expect(err).To(BeNil())

				Expect(result["id"].IsNil()).To(BeTrue())
			})
		})
	})
})
