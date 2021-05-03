package semanticid_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/happenslol/semanticid"
)

var _ = Describe("builder", func() {
	var (
		invalidSID string
		validSID   string

		testNamespace  string
		testCollection string

		testProvider semanticid.IDProvider
	)

	BeforeEach(func() {
		invalidSID = "a.b.1234"

		valid, err := semanticid.NewDefault()
		Expect(err).To(BeNil())

		validSID = valid.String()

		testNamespace = "testname"
		testCollection = "testcol"

		testProvider = &TestProvider{}
	})

	Describe("Using the semanticid builder", func() {
		Context("to build ids with default values", func() {
			It("should build the correct semanticid", func() {
				sid, err := semanticid.Builder().Build()
				Expect(err).To(BeNil())
				Expect(sid.Namespace).To(Equal(semanticid.DefaultNamespace))
				Expect(sid.Collection).To(Equal(semanticid.DefaultCollection))
				Expect(sid.ID).ToNot(BeEmpty())
			})

			It("should build valid sids from strings", func() {
				_, err := semanticid.Builder().FromString(validSID).Build()
				Expect(err).To(BeNil())
			})

			It("should reject invalid sids from strings", func() {
				_, err := semanticid.Builder().FromString(invalidSID).Build()
				Expect(err).NotTo(BeNil())
			})

			It("should accept invalid sids when validation is disabled", func() {
				_, err := semanticid.Builder().FromString(invalidSID).NoValidate().Build()
				Expect(err).To(BeNil())
			})
		})

		Context("with a custom provider", func() {
			It("should build the correct semanticid", func() {
				sid, err := semanticid.Builder().WithIDProvider(testProvider).Build()
				Expect(err).To(BeNil())
				Expect(sid.Namespace).To(Equal(semanticid.DefaultNamespace))
				Expect(sid.Collection).To(Equal(semanticid.DefaultCollection))
				Expect(sid.ID).To(Equal("1234"))
			})

			It("should accept any id value with the test provider", func() {
				_, err := semanticid.Builder().
					WithIDProvider(testProvider).
					FromString(invalidSID).
					Build()

				Expect(err).To(BeNil())
			})
		})

		Context("with custom values", func() {
			It("should build the correct semanticid", func() {
				sid, err := semanticid.Builder().
					WithNamespace(testNamespace).
					WithCollection(testCollection).
					Build()

				Expect(err).To(BeNil())
				Expect(sid.Namespace).To(Equal(testNamespace))
				Expect(sid.Collection).To(Equal(testCollection))
				Expect(sid.ID).ToNot(BeEmpty())
			})
		})
	})
})
