package semanticid_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-playground/validator/v10"
	"github.com/happenslol/semanticid"
)

type TestValidation struct {
	Both                semanticid.SemanticID `validate:"sid=test.entities"`
	OnlyNamespace       semanticid.SemanticID `validate:"sid=testname.*"`
	OnlyCollection      semanticid.SemanticID `validate:"sid=*.testcol"`
	OnlyCollectionShort semanticid.SemanticID `validate:"sid=short"`
	Any                 semanticid.SemanticID `validate:"sid=*.*"`
	Required            semanticid.SemanticID `validate:"required"`

	Pointer      *semanticid.SemanticID   `validate:"sid=testcol"`
	Array        []semanticid.SemanticID  `validate:"sid=testcol"`
	PointerArray []*semanticid.SemanticID `validate:"sid=testcol"`
}

type NoArgValidation struct {
	Test semanticid.SemanticID `validate:"sid"`
}

type InvalidArgValidation struct {
	Test semanticid.SemanticID `validate:"sid=1.2.3.4"`
}

var _ = Describe("validator", func() {
	var (
		validate   *validator.Validate
		testStruct *TestValidation
	)

	BeforeEach(func() {
		validate = validator.New()
		err := semanticid.RegisterValidation(validate)
		Expect(err).To(BeNil())

		pointerID := semanticid.Must(semanticid.New("any", "testcol"))

		testStruct = &TestValidation{
			Both:                semanticid.Must(semanticid.New("test", "entities")),
			OnlyNamespace:       semanticid.Must(semanticid.New("testname", "any")),
			OnlyCollection:      semanticid.Must(semanticid.New("any", "testcol")),
			OnlyCollectionShort: semanticid.Must(semanticid.New("any", "short")),
			Any:                 semanticid.Must(semanticid.New("01234", "56789")),
			Required:            semanticid.Must(semanticid.NewDefault()),

			Pointer: &pointerID,
			Array: []semanticid.SemanticID{
				semanticid.Must(semanticid.New("any", "testcol")),
				semanticid.Must(semanticid.New("any", "testcol")),
				semanticid.Must(semanticid.New("any", "testcol")),
				semanticid.Must(semanticid.New("any", "testcol")),
				semanticid.Must(semanticid.New("any", "testcol")),
			},
			PointerArray: []*semanticid.SemanticID{
				&pointerID,
				&pointerID,
				&pointerID,
				&pointerID,
				&pointerID,
			},
		}
	})

	Describe("Validating a struct with semanticid fields", func() {
		Context("with valid values", func() {
			It("should correctly validate", func() {
				Expect(validate.Struct(testStruct)).To(BeNil())
			})
		})

		Context("with invalid values", func() {
			It("should correctly validate namespace and collection", func() {
				testStruct.Both = semanticid.Must(semanticid.New("test", "invalid"))
				Expect(validate.Struct(testStruct)).NotTo(BeNil())

				testStruct.Both = semanticid.Must(semanticid.New("invalid", "entities"))
				Expect(validate.Struct(testStruct)).NotTo(BeNil())

				testStruct.Both = semanticid.Must(semanticid.New("invalid", "invalid"))
				Expect(validate.Struct(testStruct)).NotTo(BeNil())
			})

			It("should correctly validate namespace only", func() {
				testStruct.OnlyNamespace.Collection = "invalid"
				Expect(validate.Struct(testStruct)).To(BeNil())

				testStruct.OnlyNamespace.Namespace = "invalid"
				Expect(validate.Struct(testStruct)).NotTo(BeNil())
			})

			It("should correctly validate collection only", func() {
				testStruct.OnlyCollection.Namespace = "invalid"
				Expect(validate.Struct(testStruct)).To(BeNil())

				testStruct.OnlyCollection.Collection = "invalid"
				Expect(validate.Struct(testStruct)).NotTo(BeNil())
			})

			It("should correctly validate collection with shorthand", func() {
				testStruct.OnlyCollection.Namespace = "invalid"
				Expect(validate.Struct(testStruct)).To(BeNil())

				testStruct.OnlyCollection.Collection = "invalid"
				Expect(validate.Struct(testStruct)).NotTo(BeNil())
			})
		})

		Context("with incorrect validate tag arguments", func() {
			It("should panic when no arguments are passed", func() {
				fn := func() {
					if err := validate.Struct(&NoArgValidation{}); err != nil {
						fmt.Printf("%v\n", err)
					}
				}
				Expect(fn).To(Panic())
			})

			It("should panic when an argument with more than 1 separator is passed", func() {
				fn := func() {
					if err := validate.Struct(&InvalidArgValidation{}); err != nil {
						fmt.Printf("%v\n", err)
					}
				}
				Expect(fn).To(Panic())
			})
		})
	})
})
