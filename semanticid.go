package semanticid

import (
	"fmt"
	"strings"

	uuid "github.com/gofrs/uuid"
)

var DefaultNamespace = "namespace"
var DefaultCollection = "collection"
var Separator = "."

const (
	errUUIDError = iota
	errInvalidSID
	errInvalidUUID
	errPartContainsSeparator
)

type SemanticID struct {
	Namespace  string
	Collection string
	UUID       string
}

type SemanticIDError struct {
	errCode int
	message string
}

func (sErr SemanticIDError) Error() string {
	return sErr.message
}

func New(namespace string, collection string) (SemanticID, error) {
	uuidPart, err := uuid.NewV4()
	if err != nil {
		return SemanticID{}, SemanticIDError{
			errCode: errUUIDError,
			message: err.Error(),
		}
	}

	if strings.Contains(namespace, Separator) {
		return SemanticID{}, SemanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"namespace `%s` can't contain the separator (%s)!",
				namespace,
				Separator,
			),
		}
	}

	if strings.Contains(collection, Separator) {
		return SemanticID{}, SemanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"collection `%s` can't contain the separator (%s)!",
				collection,
				Separator,
			),
		}
	}

	return SemanticID{
		Namespace:  namespace,
		Collection: collection,
		UUID:       uuidPart.String(),
	}, nil
}

func NewWithCollection(collection string) (SemanticID, error) {
	return New(DefaultNamespace, collection)
}

func NewWithNamespace(namespace string) (SemanticID, error) {
	return New(namespace, DefaultCollection)
}

func NewDefault() (SemanticID, error) {
	return New(DefaultNamespace, DefaultCollection)
}

func FromString(s string) (SemanticID, error) {
	parts := strings.SplitN(s, Separator, 3)

	// SplitN(_, 3) guarantees at most len 3 for the
	// result, so we only need to check if there aren't enough
	if len(parts) < 3 {
		return SemanticID{}, SemanticIDError{
			errCode: errInvalidSID,
			message: fmt.Sprintf("%s is not a valid semantic id", s),
		}
	}

	namespace := parts[0]
	collection := parts[1]
	uuidPart := parts[2]

	// check if the UUID part is valid
	// TODO: Do we want a way to turn this check off? The user
	// might want to use something other than a uuid for the id
	// part
	_, err := uuid.FromString(uuidPart)
	if err != nil {
		return SemanticID{}, SemanticIDError{
			errCode: errInvalidUUID,
			message: fmt.Sprintf("The UUID section for %s is invalid", s),
		}
	}

	return SemanticID{
		Namespace:  namespace,
		Collection: collection,
		UUID:       uuidPart,
	}, nil
}

func (sID SemanticID) IsNil() bool {
	return sID.Namespace == "" && sID.Collection == "" && sID.UUID == ""
}

func (sID SemanticID) String() string {
	parts := []string{sID.Namespace, sID.Collection, sID.UUID}
	return strings.Join(parts, Separator)
}

func Must(sID SemanticID, err error) SemanticID {
	if err != nil {
		panic(err)
	}

	return sID
}
