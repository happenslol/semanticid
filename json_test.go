package semanticid_test

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/happenslol/semanticid"
)

var _ = Describe("json", func() {
	Describe("Marshalling semanticids to json", func() {
		Context("with a default semanticid", func() {
			It("should work correctly", func() {
				sid, err := semanticid.NewDefault()
				Expect(err).To(BeNil())

				m, err := json.Marshal(sid)
				Expect(err).To(BeNil())

				mStr := string(m)
				expected := fmt.Sprintf("\"%s\"", sid.String())

				Expect(mStr).To(Equal(expected))
			})
		})

		Context("with a zero value semanticid", func() {
			It("should work correctly", func() {
				zeroVal := semanticid.SemanticID{}

				m, err := json.Marshal(zeroVal)
				Expect(err).To(BeNil())

				mStr := string(m)
				Expect(mStr).To(Equal("null"))
			})
		})
	})

	Describe("Unmarshalling semanticids from json", func() {
		Context("with a valid semanticid", func() {
			sid, err := semanticid.NewDefault()
			Expect(err).To(BeNil())

			jsonStr := fmt.Sprintf("\"%s\"", sid.String())
			var result semanticid.SemanticID
			err = json.Unmarshal([]byte(jsonStr), &result)

			Expect(err).To(BeNil())

			Expect(result.Namespace).To(Equal(sid.Namespace))
			Expect(result.Collection).To(Equal(sid.Collection))
			Expect(result.ID).To(Equal(sid.ID))
		})

		Context("with a semanticid that has an invalid ID", func() {
			jsonStr := "\"namespace:collection:1234\""
			var result semanticid.SemanticID
			err := json.Unmarshal([]byte(jsonStr), &result)

			Expect(err).NotTo(BeNil())
		})

		Context("with a semanticid that has an invalid format", func() {
			jsonStr := "\"123456789\""
			var result semanticid.SemanticID
			err := json.Unmarshal([]byte(jsonStr), &result)

			Expect(err).NotTo(BeNil())
		})
	})
})
