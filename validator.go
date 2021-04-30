package semanticid

import "github.com/go-playground/validator/v10"

func SemanticIDCollectionValidation(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	collection := fl.Param()

	switch value.(type) {
	case *SemanticID:
		v := value.(*SemanticID)
		return v != nil && !v.IsNil() && v.Collection == collection

	case SemanticID:
		v := value.(SemanticID)
		return !v.IsNil() && v.Collection == collection

	case []SemanticID:
		v := value.([]SemanticID)

		for _, elem := range v {
			if elem.IsNil() || elem.Collection != collection {
				return false
			}
		}

		return true
	case []*SemanticID:
		v := value.([]*SemanticID)

		for _, elem := range v {
			if elem == nil || elem.IsNil() || elem.Collection != collection {
				return false
			}
		}

		return true
	}

	return false
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterValidation("sidcol", SemanticIDCollectionValidation)
}
