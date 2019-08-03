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
	errInvalidArgs = iota
	errUUIDError
	errInvalidSID
	errInvalidUUID
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

func New(args ...string) (SemanticID, error) {
	if len(args) > 2 {
		return SemanticID{}, SemanticIDError{}
	}

	uuidPart, err := uuid.NewV4()
	if err != nil {
		return SemanticID{}, SemanticIDError{
			errCode: errUUIDError,
			message: err.Error(),
		}
	}

	uuidStr := uuidPart.String()

	if len(args) == 1 {
		return SemanticID{
			Namespace:  args[0],
			Collection: DefaultCollection,
			UUID:       uuidStr,
		}, nil
	}

	if len(args) == 2 {
		return SemanticID{
			Namespace:  args[0],
			Collection: args[1],
			UUID:       uuidStr,
		}, nil
	}

	return SemanticID{
		Namespace:  DefaultNamespace,
		Collection: DefaultCollection,
		UUID:       uuidStr,
	}, nil
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

func IsNil(sID SemanticID) bool {
	return sID.Namespace == "" && sID.Collection == "" && sID.UUID == ""
}

func Must(sID SemanticID, err error) SemanticID {
	if err != nil {
		panic(err)
	}

	return sID
}
