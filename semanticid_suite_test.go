package semanticid_test

import (
	"testing"

	"github.com/happenslol/semanticid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSemanticid(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Semanticid Suite")
}

type TestProvider struct{}

var _ semanticid.IDProvider = &TestProvider{}

func (tp *TestProvider) Generate() (string, error) {
	return "1234", nil
}

func (tp *TestProvider) Validate(id string) error {
	return nil
}
