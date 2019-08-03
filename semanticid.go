package semanticid

import (
	"fmt"
	"strings"

	uuid "github.com/gofrs/uuid"
)

// DefaultNamespace is the namespace that will be used if
// no namespace is specified when creating a SemanticID.
var DefaultNamespace = "namespace"

// DefaultCollection is the collection that will be used if
// no collection is specified when creating a SemanticID.
var DefaultCollection = "collection"

// Separator that will be used for all SemanticIDs. You should set
// this once and never change it for your application - Once you change
// it, SemanticIDs created before that point can't be parsed anymore.
// By default, this is set to `.` since this makes SemanticIDs entirely
// URL-safe.
var Separator = "."

const (
	errUUIDError = iota
	errInvalidSID
	errInvalidUUID
	errPartContainsSeparator
)

// A SemanticID is a unique identifier for an entity that consists
// of a namespace, a collection and a UUID.
type SemanticID struct {
	Namespace  string
	Collection string
	UUID       string
}

type semanticIDError struct {
	errCode int
	message string
}

func (sErr semanticIDError) Error() string {
	return sErr.message
}

// New creates a unique SemanticID with the given namespace,
// collection and the global separator (`.` by default).
func New(namespace string, collection string) (SemanticID, error) {
	uuidPart, err := uuid.NewV4()
	if err != nil {
		return SemanticID{}, semanticIDError{
			errCode: errUUIDError,
			message: err.Error(),
		}
	}

	if strings.Contains(namespace, Separator) {
		return SemanticID{}, semanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"namespace `%s` can't contain the separator (%s)!",
				namespace,
				Separator,
			),
		}
	}

	if strings.Contains(collection, Separator) {
		return SemanticID{}, semanticIDError{
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

// NewWithCollection creates a unique SemanticID with the given
// collection and the default namespace.
func NewWithCollection(collection string) (SemanticID, error) {
	return New(DefaultNamespace, collection)
}

// NewWithNamespace creates a unique SemanticID with the given
// namespace and the default collection.
func NewWithNamespace(namespace string) (SemanticID, error) {
	return New(namespace, DefaultCollection)
}

// NewDefault creates a unique SemanticID with the default namespace
// and collection.
func NewDefault() (SemanticID, error) {
	return New(DefaultNamespace, DefaultCollection)
}

// FromString attempts to parse a given string into a SemanticID.
func FromString(s string) (SemanticID, error) {
	parts := strings.SplitN(s, Separator, 3)

	// SplitN(_, 3) guarantees at most len 3 for the
	// result, so we only need to check if there aren't enough
	if len(parts) < 3 {
		return SemanticID{}, semanticIDError{
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
		return SemanticID{}, semanticIDError{
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

// IsNil checks whether or not the SemanticID has any of its part
// set to a non-null string.
func (sID SemanticID) IsNil() bool {
	return sID.Namespace == "" && sID.Collection == "" && sID.UUID == ""
}

// String outputs a string representation of the SemanticID
func (sID SemanticID) String() string {
	parts := []string{sID.Namespace, sID.Collection, sID.UUID}
	return strings.Join(parts, Separator)
}

// Must is a convenience function that converts errors into panics on functions
// that create or parse a SemanticID.
func Must(sID SemanticID, err error) SemanticID {
	if err != nil {
		panic(err)
	}

	return sID
}
