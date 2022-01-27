package semanticid

import (
	"fmt"
	"reflect"
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
var DefaultIDProvider IDProvider = NewULIDProvider()

var empty = SemanticID{}

const (
	errIDProviderError = iota
	errInvalidSID
	errInvalidID
	errPartContainsSeparator
	errEmpty
)

var (
	ErrEmpty                 = &SemanticIDError{errEmpty, ""}
	ErrIDProvider            = &SemanticIDError{errIDProviderError, ""}
	ErrInvalid               = &SemanticIDError{errInvalidSID, ""}
	ErrInvalidIDPart         = &SemanticIDError{errInvalidID, ""}
	ErrPartContainsSeparator = &SemanticIDError{errPartContainsSeparator, ""}
)

// A SemanticID is a unique identifier for an entity that consists
// of a namespace, a collection and an ID.
type SemanticID struct {
	Namespace  string
	Collection string
	ID         string
}

type SemanticIDError struct {
	errCode int
	message string
}

func (err *SemanticIDError) Error() string {
	return err.message
}

func (err *SemanticIDError) Is(target error) bool {
	t, ok := target.(*SemanticIDError)
	if !ok {
		return false
	}

	return err.errCode == t.errCode
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

// FromStrings attempts to parse a given list of strings into a
// list of SemanticIDs. An error will be returned for the first
// conversion that errors, which means that a list that returns
// an error is not guaranteed to only contain that one error.
func FromStrings(s []string) ([]SemanticID, error) {
	result := make([]SemanticID, len(s))
	for i, id := range s {
		sid, err := FromString(id)
		if err != nil {
			return nil, err
		}

		result[i] = sid
	}

	return result, nil
}

func newWithParams(namespace, collection string, idp IDProvider) (SemanticID, error) {
	id, err := idp.Generate()
	if err != nil {
		return empty, &SemanticIDError{
			errCode: errIDProviderError,
			message: err.Error(),
		}
	}

	if strings.Contains(namespace, Separator) {
		return empty, &SemanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"Namespace `%s` can't contain the separator (%s)",
				namespace,
				Separator,
			),
		}
	}

	if strings.Contains(collection, Separator) {
		return empty, &SemanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"Collection `%s` can't contain the separator (%s)",
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
	if s == "" {
		return SemanticID{}, &SemanticIDError{errEmpty, "The given string was empty"}
	}

	parts := strings.SplitN(s, Separator, 3)

	// SplitN(_, 3) guarantees at most len 3 for the
	// result, so we only need to check if there aren't enough
	if len(parts) < 3 {
		return empty, &SemanticIDError{
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
			return empty, &SemanticIDError{
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
	return sID.Namespace == "" && sID.Collection == "" && sID.ID == ""
}

// String outputs a string representation of the SemanticID
func (sID SemanticID) String() string {
	if sID.IsNil() {
		return ""
	}

	return strings.Join([]string{sID.Namespace, sID.Collection, sID.ID}, Separator)
}

// Is checks the identity of a SemanticID, given by its Namespace and Collection.
// It expects a dot-separated Namespace and Collection combination, such that
// `semanticid.New("auth", "users").Is("auth.users") == true`.
func (sID SemanticID) Is(identity string) bool {
	if sID.IsNil() {
		return false
	}

	return fmt.Sprintf("%s.%s", sID.Namespace, sID.Collection) == identity
}

// Must is a convenience function that converts errors into panics on functions
// that create or parse a SemanticID.
func Must(sID SemanticID, err error) SemanticID {
	if err != nil {
		panic(err)
	}

	return sID
}

// CollectionForModel returns the collection defined in the `sid` tag
// on the `ID` field of the passed struct.
func CollectionForModel(model interface{}) (string, error) {
	return CollectionForModelField(model, "ID")
}

// CollectionForModelField returns the collection defined in the `sid` tag
// on the given field.
func CollectionForModelField(model interface{}, field string) (string, error) {
	var t reflect.Type
	kind := reflect.ValueOf(model).Kind()
	if kind == reflect.Struct {
		t = reflect.TypeOf(model)
	} else if kind == reflect.Ptr {
		t = reflect.Indirect(reflect.ValueOf(model)).Type()
	}

	f, ok := t.FieldByName(field)
	if !ok {
		return "", fmt.Errorf("Field `%s` not found on model", field)
	}

	tag := f.Tag.Get("sid")
	if tag == "" {
		return "", fmt.Errorf("Field `%s` did not include an sid tag", field)
	}

	if strings.Contains(tag, Separator) {
		return "", &SemanticIDError{
			errCode: errPartContainsSeparator,
			message: fmt.Sprintf(
				"Collection `%s` can't contain the separator (`%s`)",
				tag,
				Separator,
			),
		}
	}

	return tag, nil
}

// NewForModel creates a unique SemanticID using the collection defined
// in the `sid` tag in the given model.
func NewForModel(model interface{}) (SemanticID, error) {
	collection, err := CollectionForModel(model)
	if err != nil {
		return empty, err
	}

	return NewWithCollection(collection)
}
