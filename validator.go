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
	raw := fl.Field().Interface()

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

	switch value := raw.(type) {
	case string:
		sid, err := FromString(value)
		if err != nil {
			return false
		}

		return validateSIDPrefix(sid, namespace, collection)
	case []string:
		for _, s := range value {
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
	switch value := raw.(type) {
	case SemanticID:
		if value.IsNil() {
			return ""
		}

		return value.String()
	case *SemanticID:
		if value.IsNil() {
			return ""
		}

		return value.String()
	case []SemanticID:
		idStrs := make([]string, len(value))
		for i, id := range value {
			if id.IsNil() {
				idStrs[i] = ""
			} else {
				idStrs[i] = id.String()
			}
		}

		return idStrs
	case []*SemanticID:
		idStrs := make([]string, len(value))
		for i, id := range value {
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

func RegisterValidation(v *validator.Validate) error {
	v.RegisterCustomTypeFunc(
		SemanticIDTypeFunc,
		SemanticID{},
		&SemanticID{},
		[]SemanticID{},
		[]*SemanticID{},
	)

	return v.RegisterValidation("sid", SemanticIDValidation)
}
