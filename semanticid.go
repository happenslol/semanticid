package semanticid

import (
	"fmt"
	"strings"
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

// DefaultIDProvider determines the provider that will be used to
// generate and validate IDs. You can either set this or use the
// builder to select the provider on an individual basis.
var DefaultIDProvider = NewULIDProvider()

var empty = SemanticID{}

const (
	errIDProviderError = iota
	errInvalidSID
	errInvalidID
	errPartContainsSeparator
)

// A SemanticID is a unique identifier for an entity that consists
// of a namespace, a collection and an ID.
type SemanticID struct {
	Namespace  string
	Collection string
	ID         string
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
func New(namespace, collection string) (SemanticID, error) {
	return newWithParams(namespace, collection, DefaultIDProvider)
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
	return fromStringWithParams(s, DefaultIDProvider, true)
}

func newWithParams(namespace, collection string, idp IDProvider) (SemanticID, error) {
	id, err := idp.Generate()
	if err != nil {
		return empty, semanticIDError{
			errCode: errIDProviderError,
			message: err.Error(),
		}
	}

	if strings.Contains(namespace, Separator) {
		return empty, semanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"namespace `%s` can't contain the separator (%s)!",
				namespace,
				Separator,
			),
		}
	}

	if strings.Contains(collection, Separator) {
		return empty, semanticIDError{
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
		ID:         id,
	}, nil
}

func fromStringWithParams(s string, idp IDProvider, validate bool) (SemanticID, error) {
	parts := strings.SplitN(s, Separator, 3)

	// SplitN(_, 3) guarantees at most len 3 for the
	// result, so we only need to check if there aren't enough
	if len(parts) < 3 {
		return empty, semanticIDError{
			errCode: errInvalidSID,
			message: fmt.Sprintf("%s is not a valid semantic id", s),
		}
	}

	namespace := parts[0]
	collection := parts[1]
	id := parts[2]

	if validate {
		// check if the ID part is valid
		err := idp.Validate(id)
		if err != nil {
			return empty, semanticIDError{
				errCode: errInvalidID,
				message: fmt.Sprintf("The UUID section for %s is invalid", s),
			}
		}
	}

	return SemanticID{
		Namespace:  namespace,
		Collection: collection,
		ID:         id,
	}, nil
}

// IsNil checks whether or not the SemanticID has any of its part
// set to a non-null string.
func (sID SemanticID) IsNil() bool {
	return sID == empty
}

// String outputs a string representation of the SemanticID
func (sID SemanticID) String() string {
	return strings.Join([]string{sID.Namespace, sID.Collection, sID.ID}, Separator)
}

// Must is a convenience function that converts errors into panics on functions
// that create or parse a SemanticID.
func Must(sID SemanticID, err error) SemanticID {
	if err != nil {
		panic(err)
	}

	return sID
}
