package semanticid

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validateSIDPrefix(sid SemanticID, namespace, collection string) bool {
	if namespace != "*" && sid.Namespace != namespace {
		return false
	}

	if collection != "*" && sid.Collection != collection {
		return false
	}

	return true
}

func SemanticIDValidation(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()

	param := fl.Param()
	namespace := "*"
	collection := "*"
	parts := strings.Split(param, Separator)

	//NOTE(happens): We panic here in case of an invalid param,
	// similar to how the baked in validators handle this case
	switch len(parts) {
	case 2:
		namespace = parts[0]
		collection = parts[1]
	case 1:
		if parts[0] == "" {
			panic("Expected an argument for sid validator")
		}

		collection = parts[0]
	default:
		panic(fmt.Sprintf("Bad sid validation argument: %s", param))
	}

	switch value.(type) {
	case string:
		sid, err := FromString(value.(string))
		if err != nil {
			return false
		}

		return validateSIDPrefix(sid, namespace, collection)
	case []string:
		for _, s := range value.([]string) {
			sid, err := FromString(s)
			if err != nil {
				return false
			}

			if !validateSIDPrefix(sid, namespace, collection) {
				return false
			}
		}

		return true
	}

	return false
}

func SemanticIDTypeFunc(field reflect.Value) interface{} {
	raw := field.Interface()
	switch raw.(type) {
	case SemanticID:
		id := raw.(SemanticID)
		if id.IsNil() {
			return ""
		}

		return id.String()
	case *SemanticID:
		id := raw.(*SemanticID)
		if id.IsNil() {
			return ""
		}

		return id.String()
	case []SemanticID:
		ids := raw.([]SemanticID)
		idStrs := make([]string, len(ids))
		for i, id := range ids {
			if id.IsNil() {
				idStrs[i] = ""
			} else {
				idStrs[i] = id.String()
			}
		}

		return idStrs
	case []*SemanticID:
		ids := raw.([]*SemanticID)
		idStrs := make([]string, len(ids))
		for i, id := range ids {
			if id.IsNil() {
				idStrs[i] = ""
			} else {
				idStrs[i] = id.String()
			}
		}

		return idStrs
	}

	return ""
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterCustomTypeFunc(
		SemanticIDTypeFunc,
		SemanticID{},
		&SemanticID{},
		[]SemanticID{},
		[]*SemanticID{},
	)

	v.RegisterValidation("sid", SemanticIDValidation)
}
