package semanticid

import (
	"crypto/rand"
	"time"

	"github.com/gofrs/uuid"
	"github.com/oklog/ulid"
)

// IDProvider represents a type that can generate and validate the ID part
// of semanticids. By default, ULIDs will be used for this, since
// they offer lots of benefits over UUIDs, but a UUID provider
// is provided and custom providers can easily be used by either
// setting the default ID Provider or using the builder.
type IDProvider interface {
	// Generate generates a random ID
	Generate() (string, error)
	// Validate validates an existing ID
	Validate(id string) error
}

type ULIDProvider struct{}

var _ IDProvider = &ULIDProvider{}

func NewULIDProvider() *ULIDProvider {
	return &ULIDProvider{}
}

func (up *ULIDProvider) Generate() (string, error) {
	t := time.Unix(1000000, 0)
	result, err := ulid.New(ulid.Timestamp(t), rand.Reader)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (up *ULIDProvider) Validate(id string) error {
	_, err := ulid.Parse(id)
	return err
}

type UUIDProvider struct{}

var _ IDProvider = &UUIDProvider{}

func NewUUIDProvider() *UUIDProvider {
	return &UUIDProvider{}
}

func (up *UUIDProvider) Generate() (string, error) {
	result, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (up *UUIDProvider) Validate(id string) error {
	_, err := uuid.FromString(id)
	return err
}
