package semanticid_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/happenslol/semanticid"
)

type TestModel struct {
	ID         semanticid.SemanticID `sid:"testmodels"`
	CustomID   string                `sid:"custom"`
	InvalidTag string                `sid:"test.models"`
}

var _ = Describe("semanticid", func() {
	Describe("Creating semanticids", func() {
		Context("with default values", func() {
			It("should succeed", func() {
				sid, err := semanticid.NewDefault()
				Expect(err).To(BeNil())

				Expect(sid.Namespace).To(Equal("namespace"))
				Expect(sid.Collection).To(Equal("collection"))
				Expect(sid.ID).ToNot(BeEmpty())
			})
		})

		Context("with a custom namespace", func() {
			It("should succeed", func() {
				sid, err := semanticid.NewWithNamespace("testname")
				Expect(err).To(BeNil())

				Expect(sid.Namespace).To(Equal("testname"))
				Expect(sid.Collection).To(Equal("collection"))
				Expect(sid.ID).ToNot(BeEmpty())
			})
		})

		Context("with a custom collection", func() {
			It("should succeed", func() {
				sid, err := semanticid.NewWithCollection("testcol")
				Expect(err).To(BeNil())

				Expect(sid.Namespace).To(Equal("namespace"))
				Expect(sid.Collection).To(Equal("testcol"))
				Expect(sid.ID).ToNot(BeEmpty())
			})
		})

		Context("with a custom namespace and collection", func() {
			It("should succeed", func() {
				sid, err := semanticid.New("testname", "testcol")
				Expect(err).To(BeNil())

				Expect(sid.Namespace).To(Equal("testname"))
				Expect(sid.Collection).To(Equal("testcol"))
				Expect(sid.ID).ToNot(BeEmpty())
			})
		})

		Context("with invalid namespaces or collections", func() {
			It("should reject the separator in namespace and collection", func() {
				sid, err := semanticid.NewWithNamespace("test.test")
				Expect(err).NotTo(BeNil())
				Expect(sid.IsNil()).To(BeTrue())

				sid, err = semanticid.NewWithCollection("test.test")
				Expect(err).NotTo(BeNil())
				Expect(sid.IsNil()).To(BeTrue())
			})
		})

		Context("from strings", func() {
			It("should parse valid semantic ids", func() {
			})

			It("should reject invalid semantic ids", func() {
			})
		})
	})

	Describe("Using the collection struct tag", func() {
		Context("with the default model field", func() {
			It("should return the correct collection for a model value", func() {
				collection, err := semanticid.CollectionForModel(TestModel{})
				Expect(err).To(BeNil())
				Expect(collection).To(Equal("testmodels"))
			})

			It("should return the correct collection for a model pointer", func() {
				collection, err := semanticid.CollectionForModel(&TestModel{})
				Expect(err).To(BeNil())
				Expect(collection).To(Equal("testmodels"))
			})
		})

		Context("with a custom model field", func() {
			It("should return the correct collection", func() {
				collection, err := semanticid.CollectionForModelField(TestModel{}, "CustomID")
				Expect(err).To(BeNil())
				Expect(collection).To(Equal("custom"))
			})

			It("should return an error for non-existant fields", func() {
				_, err := semanticid.CollectionForModelField(TestModel{}, "NonExistant")
				Expect(err).NotTo(BeNil())
			})
		})

		Context("with an invalid tag", func() {
			It("should return an error", func() {
				_, err := semanticid.CollectionForModelField(TestModel{}, "InvalidTag")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
