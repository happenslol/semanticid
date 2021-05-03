package semanticid_test

import (
	"crypto/rand"

	"github.com/oklog/ulid"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/happenslol/semanticid"
)

var _ = Describe("id", func() {
	Describe("Using the ULID provider", func() {
		var ulidProvider semanticid.IDProvider

		BeforeEach(func() {
			ulidProvider = semanticid.NewULIDProvider()
		})

		Context("to generate IDs", func() {
			It("should generate valid IDs", func() {
				id, err := ulidProvider.Generate()
				Expect(err).To(BeNil())
				Expect(id).NotTo(BeEmpty())

				_, err = ulid.Parse(id)
				Expect(err).To(BeNil())
			})
		})

		Context("to validate IDs", func() {
			It("should accept valid IDs", func() {
				validULID, err := ulid.New(0, rand.Reader)
				Expect(err).To(BeNil())

				err = ulidProvider.Validate(validULID.String())
				Expect(err).To(BeNil())
			})

			It("should reject invalid IDs", func() {
				err := ulidProvider.Validate("1234")
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Using the UUID provider", func() {
		var uuidProvider semanticid.IDProvider

		BeforeEach(func() {
			uuidProvider = semanticid.NewUUIDProvider()
		})

		Context("to generate IDs", func() {
			It("should generate valid IDs", func() {
				id, err := uuidProvider.Generate()
				Expect(err).To(BeNil())
				Expect(id).NotTo(BeEmpty())

				_, err = uuid.FromString(id)
				Expect(err).To(BeNil())
			})
		})

		Context("to validate IDs", func() {
			It("should accept valid IDs", func() {
				validUUID, err := uuid.NewV4()
				Expect(err).To(BeNil())

				err = uuidProvider.Validate(validUUID.String())
				Expect(err).To(BeNil())
			})

			It("should reject invalid IDs", func() {
				err := uuidProvider.Validate("1234")
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
